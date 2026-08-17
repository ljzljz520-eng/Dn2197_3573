package audit

import (
	"sort"
	"strings"

	"gradebook/domain"
)

func ByOperator(events []domain.AuditEvent, operator string) []domain.AuditEvent {
	result := []domain.AuditEvent{}
	needle := strings.ToLower(strings.TrimSpace(operator))
	for _, event := range events {
		if strings.ToLower(event.Operator) == needle {
			result = append(result, event)
		}
	}
	return result
}

func ByAction(events []domain.AuditEvent, action string) []domain.AuditEvent {
	result := []domain.AuditEvent{}
	needle := strings.ToLower(strings.TrimSpace(action))
	for _, event := range events {
		if strings.ToLower(event.Action) == needle {
			result = append(result, event)
		}
	}
	return result
}

func Timeline(events []domain.AuditEvent) []domain.AuditEvent {
	result := append([]domain.AuditEvent(nil), events...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Sequence == result[j].Sequence {
			return result[i].ID < result[j].ID
		}
		return result[i].Sequence < result[j].Sequence
	})
	return result
}

func Operators(events []domain.AuditEvent) []string {
	seen := map[string]bool{}
	for _, event := range events {
		seen[event.Operator] = true
	}
	result := make([]string, 0, len(seen))
	for operator := range seen {
		result = append(result, operator)
	}
	sort.Strings(result)
	return result
}
