package domain

import (
	"fmt"
	"strings"
)

type Student struct {
	ID      string
	Name    string
	Program string
	Cohort  int
	Email   string
	Advisor string
	Active  bool
}

type CourseGrade struct {
	CourseCode string
	CourseName string
	Credits    int
	Score      float64
	Letter     string
}

type Gradebook struct {
	RecordID          string
	Student           Student
	DegreeCourses     []CourseGrade
	Seminars          []CourseGrade
	Average           float64
	Rank              int
	ConfirmationCount int
	Version           int
}

type AuditEvent struct {
	ID       string
	RecordID string
	Operator string
	Action   string
	Detail   string
	Sequence int64
}

func NewStudent(id, name, program string, cohort int, email, advisor string) Student {
	return Student{ID: strings.TrimSpace(id), Name: strings.TrimSpace(name), Program: strings.TrimSpace(program), Cohort: cohort, Email: strings.TrimSpace(email), Advisor: strings.TrimSpace(advisor), Active: true}
}

func NewCourseGrade(code, name string, credits int, score float64, letter string) CourseGrade {
	return CourseGrade{CourseCode: strings.ToUpper(strings.TrimSpace(code)), CourseName: strings.TrimSpace(name), Credits: credits, Score: score, Letter: strings.ToUpper(strings.TrimSpace(letter))}
}

func (g Gradebook) DisplayName() string { return fmt.Sprintf("%s (%s)", g.Student.Name, g.Student.ID) }

func (g Gradebook) AllCourses() []CourseGrade {
	items := make([]CourseGrade, 0, len(g.DegreeCourses)+len(g.Seminars))
	items = append(items, g.DegreeCourses...)
	items = append(items, g.Seminars...)
	return items
}

func (g Gradebook) IsComplete() bool { return len(g.DegreeCourses) == 3 }
