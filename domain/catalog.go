package domain

import (
	"sort"
	"strings"
)

type Course struct {
	Code           string
	Name           string
	Credits        int
	DegreeRequired bool
	Department     string
}

type Catalog struct{ Courses map[string]Course }

func NewCatalog(courses ...Course) Catalog {
	result := Catalog{Courses: map[string]Course{}}
	for _, course := range courses {
		result.Add(course)
	}
	return result
}

func (c *Catalog) Add(course Course) {
	if c.Courses == nil {
		c.Courses = map[string]Course{}
	}
	course.Code = strings.ToUpper(strings.TrimSpace(course.Code))
	course.Name = strings.TrimSpace(course.Name)
	c.Courses[course.Code] = course
}

func (c Catalog) Lookup(code string) (Course, bool) {
	course, ok := c.Courses[strings.ToUpper(strings.TrimSpace(code))]
	return course, ok
}

func (c Catalog) Required() []Course {
	result := []Course{}
	for _, course := range c.Courses {
		if course.DegreeRequired {
			result = append(result, course)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Code < result[j].Code })
	return result
}

func (c Catalog) DepartmentCourses(department string) []Course {
	result := []Course{}
	needle := strings.ToLower(strings.TrimSpace(department))
	for _, course := range c.Courses {
		if strings.ToLower(course.Department) == needle {
			result = append(result, course)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Code < result[j].Code })
	return result
}

func (c Catalog) Credits() int {
	total := 0
	for _, course := range c.Courses {
		total += course.Credits
	}
	return total
}

func (c Catalog) RequiredCredits() int {
	total := 0
	for _, course := range c.Courses {
		if course.DegreeRequired {
			total += course.Credits
		}
	}
	return total
}

func DefaultCatalog() Catalog {
	return NewCatalog(
		Course{Code: "GRA501", Name: "Research Methods", Credits: 3, DegreeRequired: true, Department: "graduate"},
		Course{Code: "GRA502", Name: "Advanced Theory", Credits: 3, DegreeRequired: true, Department: "graduate"},
		Course{Code: "GRA503", Name: "Scholarly Practice", Credits: 2, DegreeRequired: true, Department: "graduate"},
		Course{Code: "SEM601", Name: "Research Seminar", Credits: 2, Department: "graduate"},
	)
}
