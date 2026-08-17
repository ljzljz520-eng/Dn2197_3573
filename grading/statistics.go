package grading

import (
	"math"
	"sort"

	"gradebook/domain"
)

type Statistics struct {
	Count       int
	Mean        float64
	Median      float64
	Highest     float64
	Lowest      float64
	Distinction int
	Merit       int
	Pass        int
	Review      int
}

func Summarize(records []domain.Gradebook) Statistics {
	if len(records) == 0 {
		return Statistics{}
	}
	values := make([]float64, 0, len(records))
	stats := Statistics{Count: len(records), Lowest: math.MaxFloat64}
	for _, record := range records {
		score := WeightedAverage(record)
		values = append(values, score)
		stats.Mean += score
		if score > stats.Highest {
			stats.Highest = score
		}
		if score < stats.Lowest {
			stats.Lowest = score
		}
		switch PercentageBand(score) {
		case "distinction":
			stats.Distinction++
		case "merit":
			stats.Merit++
		case "pass":
			stats.Pass++
		default:
			stats.Review++
		}
	}
	stats.Mean = math.Round(stats.Mean/float64(stats.Count)*100) / 100
	sort.Float64s(values)
	if len(values)%2 == 1 {
		stats.Median = values[len(values)/2]
	} else {
		stats.Median = (values[len(values)/2-1] + values[len(values)/2]) / 2
	}
	return stats
}

func Distribution(records []domain.Gradebook) map[string]int {
	result := map[string]int{"distinction": 0, "merit": 0, "pass": 0, "review": 0}
	for _, record := range records {
		result[PercentageBand(WeightedAverage(record))]++
	}
	return result
}

func Top(records []domain.Gradebook, limit int) []domain.Gradebook {
	if limit < 0 {
		limit = 0
	}
	ordered := SortByAverage(records)
	if limit > len(ordered) {
		limit = len(ordered)
	}
	return ordered[:limit]
}
