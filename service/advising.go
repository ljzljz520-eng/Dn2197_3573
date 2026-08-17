package service

import (
	"fmt"
	"gradebook/domain"
	"gradebook/grading"
	"sort"
	"strings"
)

type AdvisingNote struct {
	StudentID      string
	Priority       string
	Reasons        []string
	Recommendation string
}

func (s *GradebookService) AdvisingQueue(threshold float64) ([]AdvisingNote, error) {
	records, err := s.store.List()
	if err != nil {
		return nil, err
	}
	result := []AdvisingNote{}
	for _, record := range records {
		reasons := grading.RiskReasons(record, threshold)
		if len(reasons) == 0 {
			continue
		}
		priority := "normal"
		if record.Average < 60 {
			priority = "urgent"
		} else if record.Average < threshold {
			priority = "high"
		}
		result = append(result, AdvisingNote{StudentID: record.Student.ID, Priority: priority, Reasons: reasons, Recommendation: recommend(record, reasons)})
	}
	sort.SliceStable(result, func(i, j int) bool {
		left, right := priorityValue(result[i].Priority), priorityValue(result[j].Priority)
		if left == right {
			return result[i].StudentID < result[j].StudentID
		}
		return left > right
	})
	return result, nil
}

func recommend(record domain.Gradebook, reasons []string) string {
	if record.Average < 60 {
		return "schedule academic review"
	}
	if !record.IsComplete() {
		return "complete missing degree courses"
	}
	if len(record.Seminars) == 0 {
		return "enroll in a graduate seminar"
	}
	if len(reasons) > 0 {
		return "meet with academic advisor"
	}
	return "continue current plan"
}

func priorityValue(priority string) int {
	switch priority {
	case "urgent":
		return 3
	case "high":
		return 2
	default:
		return 1
	}
}

func (note AdvisingNote) String() string {
	return fmt.Sprintf("%s [%s]: %s (%s)", note.StudentID, note.Priority, note.Recommendation, strings.Join(note.Reasons, ", "))
}

func AdvisingStudentIDs(notes []AdvisingNote) []string {
	result := make([]string, 0, len(notes))
	for _, note := range notes {
		result = append(result, note.StudentID)
	}
	return result
}
