package storage

import (
	"encoding/json"
	"errors"
	"gradebook/domain"
	"sort"
)

type Snapshot struct {
	Schema   int
	Records  []domain.Gradebook
	Audits   map[string][]domain.AuditEvent
	Metadata map[string]string
}

func (s *Store) Snapshot() (Snapshot, error) {
	records, err := s.List()
	if err != nil {
		return Snapshot{}, err
	}
	snapshot := Snapshot{Schema: 1, Records: records, Audits: map[string][]domain.AuditEvent{}, Metadata: map[string]string{}}
	for _, record := range records {
		events, e := s.Audits(record.RecordID)
		if e != nil {
			return Snapshot{}, e
		}
		snapshot.Audits[record.RecordID] = events
	}
	return snapshot, nil
}

func EncodeSnapshot(snapshot Snapshot) ([]byte, error) { return json.MarshalIndent(snapshot, "", "  ") }

func DecodeSnapshot(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return snapshot, err
	}
	if snapshot.Schema != 1 {
		return snapshot, errors.New("unsupported snapshot schema")
	}
	if snapshot.Audits == nil {
		snapshot.Audits = map[string][]domain.AuditEvent{}
	}
	if snapshot.Metadata == nil {
		snapshot.Metadata = map[string]string{}
	}
	return snapshot, nil
}

func (s *Store) Restore(snapshot Snapshot) error {
	if snapshot.Schema != 1 {
		return errors.New("unsupported snapshot schema")
	}
	for _, record := range snapshot.Records {
		if err := s.Save(record); err != nil {
			return err
		}
	}
	ids := make([]string, 0, len(snapshot.Audits))
	for id := range snapshot.Audits {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		for _, event := range snapshot.Audits[id] {
			if err := s.AppendAudit(event); err != nil {
				return err
			}
		}
	}
	for key, value := range snapshot.Metadata {
		if err := s.SetMeta(key, value); err != nil {
			return err
		}
	}
	return nil
}

func SnapshotRecordIDs(snapshot Snapshot) []string {
	result := make([]string, 0, len(snapshot.Records))
	for _, record := range snapshot.Records {
		result = append(result, record.RecordID)
	}
	sort.Strings(result)
	return result
}

func SnapshotAuditCount(snapshot Snapshot) int {
	count := 0
	for _, events := range snapshot.Audits {
		count += len(events)
	}
	return count
}
