package audit

import (
	"encoding/json"
	"fmt"
	"gradebook/domain"
	"sort"
	"strings"
)

type ActionCount struct {
	Action string
	Count  int
}

func CountActions(events []domain.AuditEvent) []ActionCount {
	counts := map[string]int{}
	for _, event := range events {
		counts[event.Action]++
	}
	result := make([]ActionCount, 0, len(counts))
	for action, count := range counts {
		result = append(result, ActionCount{Action: action, Count: count})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Action < result[j].Action })
	return result
}

func ExportJSON(events []domain.AuditEvent) ([]byte, error) {
	return json.MarshalIndent(Timeline(events), "", "  ")
}

func ExportCSV(events []domain.AuditEvent) string {
	lines := []string{"sequence,record_id,operator,action,detail"}
	for _, event := range Timeline(events) {
		lines = append(lines, fmt.Sprintf("%d,%s,%s,%s,%s", event.Sequence, csv(event.RecordID), csv(event.Operator), csv(event.Action), csv(event.Detail)))
	}
	return strings.Join(lines, "\n")
}

func csv(value string) string { return `"` + strings.ReplaceAll(value, `"`, `""`) + `"` }

func FilterSequence(events []domain.AuditEvent, first, last int64) []domain.AuditEvent {
	result := []domain.AuditEvent{}
	for _, event := range events {
		if event.Sequence >= first && event.Sequence <= last {
			result = append(result, event)
		}
	}
	return Timeline(result)
}

func Latest(events []domain.AuditEvent) (domain.AuditEvent, bool) {
	ordered := Timeline(events)
	if len(ordered) == 0 {
		return domain.AuditEvent{}, false
	}
	return ordered[len(ordered)-1], true
}

func HasAction(events []domain.AuditEvent, action string) bool {
	for _, event := range events {
		if strings.EqualFold(event.Action, action) {
			return true
		}
	}
	return false
}
