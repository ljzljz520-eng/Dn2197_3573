package ranking

import (
	"gradebook/domain"
	"gradebook/grading"
	"sort"
	"strings"
)

type CohortSummary struct {
	Cohort     int
	Count      int
	Average    float64
	TopStudent string
}

func SummarizeCohorts(records []domain.Gradebook) []CohortSummary {
	groups := map[int][]domain.Gradebook{}
	for _, record := range records {
		groups[record.Student.Cohort] = append(groups[record.Student.Cohort], record)
	}
	result := make([]CohortSummary, 0, len(groups))
	for cohort, group := range groups {
		ordered := grading.SortByAverage(group)
		stats := grading.Summarize(group)
		name := ""
		if len(ordered) > 0 {
			name = ordered[0].Student.Name
		}
		result = append(result, CohortSummary{Cohort: cohort, Count: len(group), Average: stats.Mean, TopStudent: name})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Cohort < result[j].Cohort })
	return result
}

func Search(records []domain.Gradebook, query string) []domain.Gradebook {
	needle := strings.ToLower(strings.TrimSpace(query))
	result := []domain.Gradebook{}
	for _, record := range records {
		if needle == "" || strings.Contains(strings.ToLower(record.Student.ID), needle) || strings.Contains(strings.ToLower(record.Student.Name), needle) || strings.Contains(strings.ToLower(record.Student.Program), needle) {
			result = append(result, record)
		}
	}
	return result
}

func MedianRank(records []domain.Gradebook) float64 {
	ordered := Rank(records)
	if len(ordered) == 0 {
		return 0
	}
	if len(ordered)%2 == 1 {
		return ordered[len(ordered)/2].Average
	}
	return (ordered[len(ordered)/2-1].Average + ordered[len(ordered)/2].Average) / 2
}

func ByStanding(records []domain.Gradebook, standing string) []domain.Gradebook {
	result := []domain.Gradebook{}
	for _, record := range records {
		if grading.Standing(record.Average, record.IsComplete()) == standing {
			result = append(result, record)
		}
	}
	return result
}

func RankDelta(before, after []domain.Gradebook) map[string]int {
	oldRanks := map[string]int{}
	for _, record := range Rank(before) {
		oldRanks[record.RecordID] = record.Rank
	}
	result := map[string]int{}
	for _, record := range Rank(after) {
		result[record.RecordID] = oldRanks[record.RecordID] - record.Rank
	}
	return result
}
