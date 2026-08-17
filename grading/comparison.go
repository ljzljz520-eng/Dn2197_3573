package grading

import (
	"gradebook/domain"
	"math"
	"sort"
)

type CourseComparison struct {
	CourseCode string
	Before     float64
	After      float64
	Delta      float64
}

type RecordComparison struct {
	StudentID     string
	BeforeAverage float64
	AfterAverage  float64
	AverageDelta  float64
	Courses       []CourseComparison
}

func Compare(before, after domain.Gradebook) RecordComparison {
	result := RecordComparison{StudentID: after.Student.ID, BeforeAverage: WeightedAverage(before), AfterAverage: WeightedAverage(after)}
	result.AverageDelta = round2(result.AfterAverage - result.BeforeAverage)
	beforeScores := courseScores(before)
	afterScores := courseScores(after)
	seen := map[string]bool{}
	for code, previous := range beforeScores {
		current := afterScores[code]
		result.Courses = append(result.Courses, CourseComparison{CourseCode: code, Before: previous, After: current, Delta: round2(current - previous)})
		seen[code] = true
	}
	for code, current := range afterScores {
		if !seen[code] {
			result.Courses = append(result.Courses, CourseComparison{CourseCode: code, After: current, Delta: current})
		}
	}
	sort.Slice(result.Courses, func(i, j int) bool { return result.Courses[i].CourseCode < result.Courses[j].CourseCode })
	return result
}

func courseScores(record domain.Gradebook) map[string]float64 {
	result := map[string]float64{}
	for _, course := range record.AllCourses() {
		if course.Letter != "" {
			result[course.CourseCode] = LetterScore(course.Letter)
		} else {
			result[course.CourseCode] = course.Score
		}
	}
	return result
}

func round2(value float64) float64 { return math.Round(value*100) / 100 }

func Improved(comparison RecordComparison) bool { return comparison.AverageDelta > 0 }

func DeclinedCourses(comparison RecordComparison) []CourseComparison {
	result := []CourseComparison{}
	for _, course := range comparison.Courses {
		if course.Delta < 0 {
			result = append(result, course)
		}
	}
	return result
}

func BiggestChange(comparison RecordComparison) (CourseComparison, bool) {
	if len(comparison.Courses) == 0 {
		return CourseComparison{}, false
	}
	result := comparison.Courses[0]
	for _, course := range comparison.Courses[1:] {
		if math.Abs(course.Delta) > math.Abs(result.Delta) {
			result = course
		}
	}
	return result, true
}

func Trend(records []domain.Gradebook) string {
	if len(records) < 2 {
		return "insufficient data"
	}
	first := WeightedAverage(records[0])
	last := WeightedAverage(records[len(records)-1])
	if last > first {
		return "improving"
	}
	if last < first {
		return "declining"
	}
	return "stable"
}
