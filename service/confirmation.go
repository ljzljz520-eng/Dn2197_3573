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
	if gate != nil {
		// Stage a read so the barrier can coordinate both writers, then wait
		// until both operators are released before mutating the record.
		gate.Reads <- operator
		<-gate.Release
	}
	// ConfirmRecord reads and writes the record inside a single transaction, so
	// a stale read staged above does not clobber a concurrent confirmation and
	// each operator's audit event persists under a distinct key.
	_, err := s.store.ConfirmRecord(recordID, operator)
	return err
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
