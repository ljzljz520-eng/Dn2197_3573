package audit

import (
	"fmt"
	"gradebook/domain"
	"sort"
	"strings"
)

type ComplianceIssue struct {
	RecordID string
	Code     string
	Detail   string
}

type ComplianceReport struct {
	Records   int
	Events    int
	Issues    []ComplianceIssue
	Operators []string
}

func Review(records []domain.Gradebook, events map[string][]domain.AuditEvent) ComplianceReport {
	report := ComplianceReport{Records: len(records)}
	operatorSet := map[string]bool{}
	for _, record := range records {
		recordEvents := events[record.RecordID]
		report.Events += len(recordEvents)
		if len(recordEvents) == 0 {
			report.Issues = append(report.Issues, ComplianceIssue{RecordID: record.RecordID, Code: "missing_audit", Detail: "record has no audit trail"})
		}
		if !HasAction(recordEvents, "create") {
			report.Issues = append(report.Issues, ComplianceIssue{RecordID: record.RecordID, Code: "missing_creation", Detail: "creation event is absent"})
		}
		confirmations := len(ByAction(recordEvents, "confirm"))
		if confirmations != record.ConfirmationCount {
			report.Issues = append(report.Issues, ComplianceIssue{RecordID: record.RecordID, Code: "confirmation_mismatch", Detail: fmt.Sprintf("record=%d audit=%d", record.ConfirmationCount, confirmations)})
		}
		for _, event := range recordEvents {
			if strings.TrimSpace(event.Operator) == "" {
				report.Issues = append(report.Issues, ComplianceIssue{RecordID: record.RecordID, Code: "missing_operator", Detail: event.ID})
			} else {
				operatorSet[event.Operator] = true
			}
		}
	}
	for operator := range operatorSet {
		report.Operators = append(report.Operators, operator)
	}
	sort.Strings(report.Operators)
	sort.SliceStable(report.Issues, func(i, j int) bool {
		if report.Issues[i].RecordID == report.Issues[j].RecordID {
			return report.Issues[i].Code < report.Issues[j].Code
		}
		return report.Issues[i].RecordID < report.Issues[j].RecordID
	})
	return report
}

func (r ComplianceReport) Passed() bool { return len(r.Issues) == 0 }

func (r ComplianceReport) IssuesFor(recordID string) []ComplianceIssue {
	result := []ComplianceIssue{}
	for _, issue := range r.Issues {
		if issue.RecordID == recordID {
			result = append(result, issue)
		}
	}
	return result
}

func (r ComplianceReport) Summary() string {
	return fmt.Sprintf("records=%d events=%d issues=%d operators=%d", r.Records, r.Events, len(r.Issues), len(r.Operators))
}

func IssueCodes(report ComplianceReport) []string {
	set := map[string]bool{}
	for _, issue := range report.Issues {
		set[issue.Code] = true
	}
	result := make([]string, 0, len(set))
	for code := range set {
		result = append(result, code)
	}
	sort.Strings(result)
	return result
}
