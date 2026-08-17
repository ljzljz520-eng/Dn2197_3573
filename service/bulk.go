package service

import (
	"fmt"
	"gradebook/audit"
	"gradebook/domain"
	"gradebook/grading"
	"sort"
	"strings"
)

type ImportResult struct {
	Added    int
	Rejected int
	Errors   []string
}

func (s *GradebookService) AddMany(records []domain.Gradebook, operator string) ImportResult {
	result := ImportResult{}
	for _, record := range records {
		if err := s.Add(record, operator); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Added++
		}
	}
	return result
}

func (s *GradebookService) UpdateMany(records []domain.Gradebook, operator string) ImportResult {
	result := ImportResult{}
	for _, record := range records {
		if err := s.Update(record, operator, "batch update"); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Added++
		}
	}
	return result
}

func (s *GradebookService) AuditForOperator(operator string) ([]domain.AuditEvent, error) {
	records, err := s.store.List()
	if err != nil {
		return nil, err
	}
	result := []domain.AuditEvent{}
	for _, record := range records {
		events, e := s.store.Audits(record.RecordID)
		if e != nil {
			return nil, e
		}
		result = append(result, audit.ByOperator(events, operator)...)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Sequence < result[j].Sequence })
	return result, nil
}

func (s *GradebookService) CourseAverages() (map[string]float64, error) {
	records, err := s.store.List()
	if err != nil {
		return nil, err
	}
	totals := map[string]float64{}
	counts := map[string]int{}
	for _, record := range records {
		for _, course := range record.AllCourses() {
			score := course.Score
			if course.Letter != "" {
				score = grading.LetterScore(course.Letter)
			}
			totals[course.CourseCode] += score
			counts[course.CourseCode]++
		}
	}
	for code, total := range totals {
		totals[code] = total / float64(counts[code])
	}
	return totals, nil
}

func (s *GradebookService) ProgramSummary() (map[string]grading.Statistics, error) {
	records, err := s.store.List()
	if err != nil {
		return nil, err
	}
	groups := map[string][]domain.Gradebook{}
	for _, record := range records {
		key := strings.TrimSpace(record.Student.Program)
		groups[key] = append(groups[key], record)
	}
	result := map[string]grading.Statistics{}
	for program, group := range groups {
		result[program] = grading.Summarize(group)
	}
	return result, nil
}

func FormatImport(result ImportResult) string {
	return fmt.Sprintf("added=%d rejected=%d errors=%s", result.Added, result.Rejected, strings.Join(result.Errors, "; "))
}
