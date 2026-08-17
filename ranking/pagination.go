package ranking

import (
	"fmt"
	"gradebook/domain"
	"gradebook/grading"
)

type Page struct {
	Items       []domain.Gradebook
	Number      int
	Size        int
	TotalItems  int
	TotalPages  int
	HasPrevious bool
	HasNext     bool
}

func Paginate(records []domain.Gradebook, number, size int) Page {
	if size <= 0 {
		size = 20
	}
	if number <= 0 {
		number = 1
	}
	ordered := grading.SortByAverage(records)
	totalPages := len(ordered) / size
	if len(ordered)%size != 0 {
		totalPages++
	}
	start := (number - 1) * size
	if start > len(ordered) {
		start = len(ordered)
	}
	end := start + size
	if end > len(ordered) {
		end = len(ordered)
	}
	return Page{Items: append([]domain.Gradebook(nil), ordered[start:end]...), Number: number, Size: size, TotalItems: len(ordered), TotalPages: totalPages, HasPrevious: number > 1, HasNext: number < totalPages}
}

func (p Page) Empty() bool { return len(p.Items) == 0 }

func (p Page) Range() string {
	if p.Empty() {
		return "0-0"
	}
	start := (p.Number-1)*p.Size + 1
	end := start + len(p.Items) - 1
	return fmt.Sprintf("%d-%d", start, end)
}

func PageNumbers(page Page, radius int) []int {
	if radius < 0 {
		radius = 0
	}
	first := page.Number - radius
	if first < 1 {
		first = 1
	}
	last := page.Number + radius
	if last > page.TotalPages {
		last = page.TotalPages
	}
	result := []int{}
	for number := first; number <= last; number++ {
		result = append(result, number)
	}
	return result
}
