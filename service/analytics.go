package service

import (
	"gradebook/domain"
	"gradebook/grading"
	"sort"
	"strings"
)

type Analytics struct {
	Records        int
	Average        float64
	Median         float64
	AtRisk         []string
	Top            []domain.Gradebook
	CourseAverages map[string]float64
}

func (s *GradebookService) Analyze(threshold float64, top int) (Analytics, error) {
	records, err := s.store.List()
	if err != nil {
		return Analytics{}, err
	}
	stats := grading.Summarize(records)
	result := Analytics{Records: len(records), Average: stats.Mean, Median: stats.Median, Top: grading.Top(records, top), CourseAverages: map[string]float64{}}
	for _, record := range records {
		if grading.AtRisk(record, threshold) {
			result.AtRisk = append(result.AtRisk, record.Student.ID)
		}
	}
	courseAverages, err := s.CourseAverages()
	if err != nil {
		return Analytics{}, err
	}
	result.CourseAverages = courseAverages
	sort.Strings(result.AtRisk)
	return result, nil
}

func (a Analytics) HasRisk(studentID string) bool {
	for _, id := range a.AtRisk {
		if strings.EqualFold(id, studentID) {
			return true
		}
	}
	return false
}

func (a Analytics) TopNames() []string {
	names := make([]string, 0, len(a.Top))
	for _, record := range a.Top {
		names = append(names, record.Student.Name)
	}
	return names
}

func (a Analytics) CourseAverage(code string) float64 {
	return a.CourseAverages[strings.ToUpper(strings.TrimSpace(code))]
}
