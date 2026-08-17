package service

import (
	"fmt"
	"gradebook/domain"
	"gradebook/grading"
	"sort"
	"strings"
)

type TranscriptQuery struct {
	StudentID       string
	IncludeSeminars bool
	MinimumScore    float64
	CoursePrefix    string
}

func (s *GradebookService) Transcript(recordID string) (grading.Transcript, error) {
	record, err := s.store.Get(recordID)
	if err != nil {
		return grading.Transcript{}, err
	}
	return grading.BuildTranscript(record), nil
}

func (s *GradebookService) QueryTranscript(query TranscriptQuery) (grading.Transcript, error) {
	record, err := s.store.Get(strings.TrimSpace(query.StudentID))
	if err != nil {
		return grading.Transcript{}, err
	}
	filtered := record
	filtered.DegreeCourses = filterCourses(record.DegreeCourses, query.MinimumScore, query.CoursePrefix)
	if query.IncludeSeminars {
		filtered.Seminars = filterCourses(record.Seminars, query.MinimumScore, query.CoursePrefix)
	} else {
		filtered.Seminars = nil
	}
	return grading.BuildTranscript(filtered), nil
}

func filterCourses(courses []domain.CourseGrade, minimum float64, prefix string) []domain.CourseGrade {
	result := []domain.CourseGrade{}
	needle := strings.ToUpper(strings.TrimSpace(prefix))
	for _, course := range courses {
		score := course.Score
		if course.Letter != "" {
			score = grading.LetterScore(course.Letter)
		}
		if score >= minimum && (needle == "" || strings.HasPrefix(course.CourseCode, needle)) {
			result = append(result, course)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CourseCode < result[j].CourseCode })
	return result
}

func (s *GradebookService) GraduationReview(recordID string, policy domain.GradePolicy) ([]domain.PolicyViolation, error) {
	record, err := s.store.Get(recordID)
	if err != nil {
		return nil, err
	}
	return domain.EvaluatePolicy(record, policy), nil
}

func (s *GradebookService) GraduationEligible(recordID string, policy domain.GradePolicy) (bool, error) {
	violations, err := s.GraduationReview(recordID, policy)
	if err != nil {
		return false, err
	}
	return len(violations) == 0, nil
}

func (s *GradebookService) TranscriptText(recordID string) (string, error) {
	transcript, err := s.Transcript(recordID)
	if err != nil {
		return "", err
	}
	return grading.RenderTranscript(transcript), nil
}

func FormatViolations(violations []domain.PolicyViolation) string {
	if len(violations) == 0 {
		return "eligible"
	}
	return fmt.Sprintf("ineligible: %s", strings.Join(domain.ViolationMessages(violations), "; "))
}
