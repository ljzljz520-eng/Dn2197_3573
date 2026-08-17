package grading

import (
	"math"
	"sort"

	"gradebook/domain"
)

func WeightedAverage(record domain.Gradebook) float64 {
	var points float64
	var credits int
	for _, course := range record.DegreeCourses {
		points += course.Score * float64(course.Credits)
		credits += course.Credits
	}
	for _, course := range record.Seminars {
		if course.Letter != "" {
			points += LetterScore(course.Letter) * float64(course.Credits)
		} else {
			points += course.Score * float64(course.Credits)
		}
		credits += course.Credits
	}
	if credits == 0 {
		return 0
	}
	return math.Round(points/float64(credits)*100) / 100
}

func LetterScore(letter string) float64 {
	switch letter {
	case "A+":
		return 100
	case "A":
		return 95
	case "A-":
		return 90
	case "B+":
		return 88
	case "B":
		return 85
	case "B-":
		return 80
	case "C+":
		return 78
	case "C":
		return 75
	default:
		return 0
	}
}

func PercentageBand(score float64) string {
	switch {
	case score >= 90:
		return "distinction"
	case score >= 80:
		return "merit"
	case score >= 60:
		return "pass"
	default:
		return "review"
	}
}

func SortByAverage(records []domain.Gradebook) []domain.Gradebook {
	result := append([]domain.Gradebook(nil), records...)
	for i := range result {
		result[i].Average = WeightedAverage(result[i])
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Average == result[j].Average {
			return result[i].Student.ID < result[j].Student.ID
		}
		return result[i].Average > result[j].Average
	})
	for i := range result {
		result[i].Rank = i + 1
	}
	return result
}
