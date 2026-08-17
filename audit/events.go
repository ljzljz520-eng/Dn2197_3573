package audit

import (
	"fmt"
	"sync/atomic"
	"time"

	"gradebook/domain"
)

type Clock interface{ Now() time.Time }

type FixedClock struct{ value time.Time }

func NewFixedClock(value time.Time) FixedClock { return FixedClock{value: value} }
func (c FixedClock) Now() time.Time            { return c.value }

type Sequencer struct{ next atomic.Int64 }

func (s *Sequencer) Next() int64 { return s.next.Add(1) }

func NewConfirmationEvent(seq int64, recordID, operator string) domain.AuditEvent {
	return domain.AuditEvent{ID: fmt.Sprintf("confirm-%s-%d", operator, seq), RecordID: recordID, Operator: operator, Action: "confirm", Detail: "operator confirmed gradebook", Sequence: seq}
}

func NewCreationEvent(seq int64, recordID, operator string) domain.AuditEvent {
	return domain.AuditEvent{ID: fmt.Sprintf("create-%s-%d", operator, seq), RecordID: recordID, Operator: operator, Action: "create", Detail: "gradebook created", Sequence: seq}
}

func NewUpdateEvent(seq int64, recordID, operator, detail string) domain.AuditEvent {
	return domain.AuditEvent{ID: fmt.Sprintf("update-%s-%d", operator, seq), RecordID: recordID, Operator: operator, Action: "update", Detail: detail, Sequence: seq}
}
