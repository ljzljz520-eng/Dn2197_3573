package audit

import (
	"gradebook/domain"
	"testing"
	"time"
)

func TestAuditQueries(t *testing.T) {
	clock := NewFixedClock(time.Unix(100, 0))
	if clock.Now().Unix() != 100 {
		t.Fatal("clock mismatch")
	}
	events := []domain.AuditEvent{NewUpdateEvent(2, "S1", "bob", "changed"), NewConfirmationEvent(1, "S1", "ann"), NewConfirmationEvent(3, "S1", "ann")}
	if len(ByOperator(events, "ann")) != 2 || len(ByAction(events, "confirm")) != 2 || len(Timeline(events)) != 3 {
		t.Fatal("audit query mismatch")
	}
	if got := Operators(events); len(got) != 2 || got[0] != "ann" {
		t.Fatalf("operators %#v", got)
	}
}
