package service

import (
	"fmt"
	"strings"

	"gradebook/audit"
	"gradebook/domain"
	"gradebook/grading"
)

type Report struct {
	Title      string
	Lines      []string
	Statistics grading.Statistics
	Audits     []domain.AuditEvent
}

func (s *GradebookService) BuildReport(title string) (Report, error) {
	records, err := s.store.List()
	if err != nil {
		return Report{}, err
	}
	report := Report{Title: strings.TrimSpace(title), Statistics: grading.Summarize(records), Lines: []string{}}
	for _, record := range grading.SortByAverage(records) {
		report.Lines = append(report.Lines, fmt.Sprintf("%d. %s %.2f", record.Rank, record.Student.Name, record.Average))
	}
	for _, record := range records {
		events, e := s.store.Audits(record.RecordID)
		if e != nil {
			return Report{}, e
		}
		report.Audits = append(report.Audits, audit.Timeline(events)...)
	}
	return report, nil
}

func RenderReport(report Report) string {
	lines := append([]string{report.Title}, report.Lines...)
	lines = append(lines, fmt.Sprintf("mean=%.2f median=%.2f count=%d", report.Statistics.Mean, report.Statistics.Median, report.Statistics.Count))
	return strings.Join(lines, "\n")
}
