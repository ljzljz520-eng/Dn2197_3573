package grading

import (
	"fmt"
	"gradebook/domain"
	"sort"
)

type Threshold struct {
	Name        string
	Minimum     float64
	Maximum     float64
	Description string
}

func DefaultThresholds() []Threshold {
	return []Threshold{{"distinction", 90, 100, "exceptional mastery"}, {"merit", 80, 89.99, "strong mastery"}, {"pass", 60, 79.99, "satisfactory mastery"}, {"review", 0, 59.99, "requires support"}}
}

func ThresholdFor(score float64, thresholds []Threshold) Threshold {
	ordered := append([]Threshold(nil), thresholds...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Minimum > ordered[j].Minimum })
	for _, threshold := range ordered {
		if score >= threshold.Minimum && score <= threshold.Maximum {
			return threshold
		}
	}
	return Threshold{Name: "unknown", Minimum: score, Maximum: score}
}

func ValidateThresholds(thresholds []Threshold) error {
	if len(thresholds) == 0 {
		return fmt.Errorf("thresholds are empty")
	}
	ordered := append([]Threshold(nil), thresholds...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Minimum < ordered[j].Minimum })
	if ordered[0].Minimum != 0 {
		return fmt.Errorf("thresholds must start at zero")
	}
	for i := 1; i < len(ordered); i++ {
		if ordered[i].Minimum < ordered[i-1].Minimum {
			return fmt.Errorf("thresholds overlap")
		}
	}
	return nil
}

func AtRisk(record domain.Gradebook, threshold float64) bool {
	return WeightedAverage(record) < threshold || !record.IsComplete()
}

func RiskReasons(record domain.Gradebook, threshold float64) []string {
	reasons := []string{}
	if !record.IsComplete() {
		reasons = append(reasons, "degree courses incomplete")
	}
	if WeightedAverage(record) < threshold {
		reasons = append(reasons, "average below threshold")
	}
	if len(record.Seminars) == 0 {
		reasons = append(reasons, "no seminar recorded")
	}
	return reasons
}

func GradePoint(score float64) float64 {
	switch {
	case score >= 90:
		return 4
	case score >= 80:
		return 3
	case score >= 70:
		return 2
	case score >= 60:
		return 1
	default:
		return 0
	}
}

func GPA(record domain.Gradebook) float64 {
	var total float64
	var credits int
	for _, course := range record.AllCourses() {
		score := course.Score
		if course.Letter != "" {
			score = LetterScore(course.Letter)
		}
		total += GradePoint(score) * float64(course.Credits)
		credits += course.Credits
	}
	if credits == 0 {
		return 0
	}
	return total / float64(credits)
}
