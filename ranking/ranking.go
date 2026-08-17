package ranking

import (
	"gradebook/domain"
	"sort"
	"strings"
)

func Rank(records []domain.Gradebook) []domain.Gradebook {
	result := append([]domain.Gradebook(nil), records...)
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

func TopByRank(records []domain.Gradebook, count int) []domain.Gradebook {
	if count < 0 {
		count = 0
	}
	if count > len(records) {
		count = len(records)
	}
	return append([]domain.Gradebook(nil), records[:count]...)
}

func FindByCohort(records []domain.Gradebook, cohort int) []domain.Gradebook {
	result := []domain.Gradebook{}
	for _, record := range records {
		if record.Student.Cohort == cohort {
			result = append(result, record)
		}
	}
	return result
}

func FindByName(records []domain.Gradebook, name string) []domain.Gradebook {
	result := []domain.Gradebook{}
	needle := strings.ToLower(strings.TrimSpace(name))
	for _, record := range records {
		if strings.Contains(strings.ToLower(record.Student.Name), needle) {
			result = append(result, record)
		}
	}
	return result
}
