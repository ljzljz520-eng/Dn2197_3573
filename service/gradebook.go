package service

import (
	"errors"
	"fmt"

	"gradebook/audit"
	"gradebook/domain"
	"gradebook/grading"
	"gradebook/storage"
)

type GradebookService struct {
	store *storage.Store
	seq   *audit.Sequencer
}

func NewGradebookService(store *storage.Store) *GradebookService {
	return &GradebookService{store: store, seq: &audit.Sequencer{}}
}

func (s *GradebookService) Add(record domain.Gradebook, operator string) error {
	if s.store == nil {
		return errors.New("store is required")
	}
	domain.Normalize(&record)
	if err := domain.ValidateRecord(record); err != nil {
		return err
	}
	if record.RecordID == "" {
		record.RecordID = record.Student.ID
	}
	record.Average = grading.WeightedAverage(record)
	record.Version = 1
	if err := s.store.Save(record); err != nil {
		return err
	}
	event := audit.NewCreationEvent(s.seq.Next(), record.RecordID, operator)
	return s.store.AppendAudit(event)
}

func (s *GradebookService) Update(record domain.Gradebook, operator, detail string) error {
	if s.store == nil {
		return errors.New("store is required")
	}
	domain.Normalize(&record)
	if err := domain.ValidateRecord(record); err != nil {
		return err
	}
	current, err := s.store.Get(record.RecordID)
	if err != nil {
		return err
	}
	if record.Version != current.Version {
		return fmt.Errorf("version conflict: got %d want %d", record.Version, current.Version)
	}
	record.Average = grading.WeightedAverage(record)
	record.Version++
	if err := s.store.Save(record); err != nil {
		return err
	}
	event := audit.NewUpdateEvent(s.seq.Next(), record.RecordID, operator, detail)
	return s.store.AppendAudit(event)
}

func (s *GradebookService) Get(recordID string) (domain.Gradebook, error) {
	return s.store.Get(recordID)
}

func (s *GradebookService) List() ([]domain.Gradebook, error) { return s.store.List() }

func (s *GradebookService) Statistics() (grading.Statistics, error) {
	records, err := s.store.List()
	if err != nil {
		return grading.Statistics{}, err
	}
	return grading.Summarize(records), nil
}
