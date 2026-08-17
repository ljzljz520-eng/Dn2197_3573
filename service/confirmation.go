package service

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
	if _, err := s.store.Get(recordID); err != nil {
		return err
	}
	if gate != nil {
		gate.Reads <- operator
		<-gate.Release
	}
	return s.store.Confirm(recordID, operator)
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
