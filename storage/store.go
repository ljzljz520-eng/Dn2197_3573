package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"go.etcd.io/bbolt"
	"gradebook/domain"
)

var (
	gradesBucket = []byte("gradebooks")
	auditBucket  = []byte("audits")
	metaBucket   = []byte("metadata")
	ErrNotFound  = errors.New("gradebook not found")
)

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0600, nil)
	if err != nil {
		return nil, err
	}
	store := &Store{db: db, path: path}
	if err := db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{gradesBucket, auditBucket, metaBucket} {
			if _, e := tx.CreateBucketIfNotExists(bucket); e != nil {
				return e
			}
		}
		return nil
	}); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func (s *Store) Save(record domain.Gradebook) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(gradesBucket).Put([]byte(record.RecordID), data) })
}

func (s *Store) Get(recordID string) (domain.Gradebook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var record domain.Gradebook
	if s.db == nil {
		return record, errors.New("store closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(gradesBucket).Get([]byte(recordID))
		if value == nil {
			return ErrNotFound
		}
		return json.Unmarshal(append([]byte(nil), value...), &record)
	})
	return record, err
}

func (s *Store) List() ([]domain.Gradebook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []domain.Gradebook{}
	if s.db == nil {
		return result, errors.New("store closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(gradesBucket).ForEach(func(_, value []byte) error {
			var record domain.Gradebook
			if err := json.Unmarshal(value, &record); err != nil {
				return err
			}
			result = append(result, record)
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].RecordID < result[j].RecordID })
	return result, err
}

func (s *Store) AppendAudit(event domain.AuditEvent) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%020d:%s", event.RecordID, event.Sequence, event.ID)
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(auditBucket).Put([]byte(key), data) })
}

// ConfirmRecord atomically confirms a gradebook for an operator within a single
// transaction. Reading the current count inside the write transaction means a
// concurrent confirmation that staged a stale read before the writers were
// released still observes the count committed by the other writer, so the two
// increments compose into a final count of two rather than overwriting each
// other. The audit event key embeds both the authoritative per-record sequence
// and the operator, so the two events land under distinct keys and both survive.
func (s *Store) ConfirmRecord(recordID, operator string) (domain.Gradebook, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var record domain.Gradebook
	if s.db == nil {
		return record, errors.New("store closed")
	}
	err := s.db.Update(func(tx *bbolt.Tx) error {
		grades := tx.Bucket(gradesBucket)
		value := grades.Get([]byte(recordID))
		if value == nil {
			return ErrNotFound
		}
		if err := json.Unmarshal(append([]byte(nil), value...), &record); err != nil {
			return err
		}
		record.ConfirmationCount++
		record.Version++
		data, err := json.Marshal(record)
		if err != nil {
			return err
		}
		if err := grades.Put([]byte(recordID), data); err != nil {
			return err
		}
		sequence := int64(record.ConfirmationCount)
		event := domain.AuditEvent{
			ID:       fmt.Sprintf("confirmation-%s-%d", operator, sequence),
			RecordID: recordID,
			Operator: operator,
			Action:   "confirm",
			Detail:   "operator confirmed gradebook",
			Sequence: sequence,
		}
		auditData, err := json.Marshal(event)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s:%020d:%s", event.RecordID, event.Sequence, event.ID)
		return tx.Bucket(auditBucket).Put([]byte(key), auditData)
	})
	return record, err
}

func (s *Store) Audits(recordID string) ([]domain.AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []domain.AuditEvent{}
	if s.db == nil {
		return result, errors.New("store closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket(auditBucket).Cursor()
		prefix := []byte(recordID + ":")
		for key, value := cursor.Seek(prefix); key != nil && len(key) >= len(prefix) && string(key[:len(prefix)]) == string(prefix); key, value = cursor.Next() {
			var event domain.AuditEvent
			if err := json.Unmarshal(value, &event); err != nil {
				return err
			}
			result = append(result, event)
		}
		return nil
	})
	sort.Slice(result, func(i, j int) bool { return result[i].Sequence < result[j].Sequence })
	return result, err
}

func (s *Store) SetMeta(key, value string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(metaBucket).Put([]byte(key), []byte(value)) })
}

func (s *Store) GetMeta(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return "", errors.New("store closed")
	}
	var value string
	err := s.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(metaBucket).Get([]byte(key))
		if raw == nil {
			return ErrNotFound
		}
		value = string(raw)
		return nil
	})
	return value, err
}
