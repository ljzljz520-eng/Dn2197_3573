package storage

import (
	"sort"
	"strings"

	"gradebook/domain"
)

func FilterByProgram(records []domain.Gradebook, program string) []domain.Gradebook {
	result := []domain.Gradebook{}
	needle := strings.ToLower(strings.TrimSpace(program))
	for _, record := range records {
		if needle == "" || strings.ToLower(record.Student.Program) == needle {
			result = append(result, record)
		}
	}
	return result
}

func FilterActive(records []domain.Gradebook) []domain.Gradebook {
	result := []domain.Gradebook{}
	for _, record := range records {
		if record.Student.Active {
			result = append(result, record)
		}
	}
	return result
}

func SortByName(records []domain.Gradebook) []domain.Gradebook {
	result := append([]domain.Gradebook(nil), records...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Student.Name < result[j].Student.Name })
	return result
}
