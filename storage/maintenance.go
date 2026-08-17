package storage

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"gradebook/domain"
	"time"
)

type Health struct {
	Open       bool
	Gradebooks int
	Audits     int
	Path       string
}

func (s *Store) Health() (Health, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return Health{Path: s.path}, errors.New("store closed")
	}
	health := Health{Open: true, Path: s.path}
	err := s.db.View(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(gradesBucket).ForEach(func(_, _ []byte) error { health.Gradebooks++; return nil }); err != nil {
			return err
		}
		return tx.Bucket(auditBucket).ForEach(func(_, _ []byte) error { health.Audits++; return nil })
	})
	return health, err
}

func (s *Store) Delete(recordID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(gradesBucket).Delete([]byte(recordID)); err != nil {
			return err
		}
		return nil
	})
}

func (s *Store) Backup(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error { return tx.CopyFile(path, 0600) })
}

func (s *Store) Touch(label string) error {
	if label == "" {
		return errors.New("label required")
	}
	return s.SetMeta("touch:"+label, fmt.Sprintf("%d", time.Unix(0, 0).Unix()))
}

func (s *Store) SaveMany(records []domain.Gradebook) error {
	for _, record := range records {
		if err := s.Save(record); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) AuditCount(recordID string) (int, error) {
	events, err := s.Audits(recordID)
	return len(events), err
}
