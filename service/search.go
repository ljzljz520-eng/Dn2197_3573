package service

import (
	"sort"
	"strings"

	"gradebook/domain"
	"gradebook/grading"
	"gradebook/storage"
)

func (s *GradebookService) ByID(id string) (domain.Gradebook, error) {
	return s.store.Get(strings.TrimSpace(id))
}

func (s *GradebookService) ByName(name string) ([]domain.Gradebook, error) {
	records, err := s.store.List()
	if err != nil {
		return nil, err
	}
	needle := strings.ToLower(strings.TrimSpace(name))
	result := []domain.Gradebook{}
	for _, record := range records {
		if strings.Contains(strings.ToLower(record.Student.Name), needle) {
			result = append(result, record)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Student.Name < result[j].Student.Name })
	return result, nil
}

func (s *GradebookService) Ranked(limit int) ([]domain.Gradebook, error) {
	records, err := s.store.List()
	if err != nil {
		return nil, err
	}
	if limit == 0 {
		limit = len(records)
	}
	return grading.Top(records, limit), nil
}

func (s *GradebookService) ByProgram(program string) ([]domain.Gradebook, error) {
	records, err := s.store.List()
	if err != nil {
		return nil, err
	}
	return storage.FilterByProgram(records, program), nil
}
