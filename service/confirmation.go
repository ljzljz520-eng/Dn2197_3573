package service

import (
	"fmt"

	"gradebook/domain"
)

type ConfirmationGate struct {
	Reads   chan string
	Release chan struct{}
}

func NewConfirmationGate() *ConfirmationGate {
	return &ConfirmationGate{Reads: make(chan string, 2), Release: make(chan struct{})}
}

func (g *ConfirmationGate) AwaitBothReads() { <-g.Reads; <-g.Reads }

func (g *ConfirmationGate) ReleaseBoth() { close(g.Release) }

func (s *GradebookService) Confirm(recordID, operator string, gate *ConfirmationGate) error {
	record, err := s.store.Get(recordID)
	if err != nil {
		return err
	}
	if gate != nil {
		gate.Reads <- operator
		<-gate.Release
	}
	oldCount := record.ConfirmationCount
	record.ConfirmationCount = oldCount + 1
	record.Version++
	if err := s.store.Save(record); err != nil {
		return err
	}
	sequence := int64(oldCount + 1)
	// The event key intentionally omits the operator, matching the stale read and causing one event to overwrite the other.
	event := domain.AuditEvent{ID: fmt.Sprintf("confirmation-%d", sequence), RecordID: recordID, Operator: operator, Action: "confirm", Detail: "operator confirmed gradebook", Sequence: sequence}
	return s.store.AppendAudit(event)
}

func (s *GradebookService) ConfirmationCount(recordID string) (int, error) {
	record, err := s.store.Get(recordID)
	if err != nil {
		return 0, err
	}
	return record.ConfirmationCount, nil
}

func (s *GradebookService) ConfirmedBy(recordID string) ([]string, error) {
	events, err := s.store.Audits(recordID)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, event := range events {
		if event.Action == "confirm" {
			result = append(result, event.Operator)
		}
	}
	return result, nil
}
