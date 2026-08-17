package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	ErrInvalidStudent  = errors.New("invalid student")
	ErrInvalidCourse   = errors.New("invalid course")
	ErrDuplicateCourse = errors.New("duplicate course")
)

func ValidateStudent(s Student) error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("%w: id and name are required", ErrInvalidStudent)
	}
	if s.Cohort < 2000 || s.Cohort > 2100 {
		return fmt.Errorf("%w: cohort out of range", ErrInvalidStudent)
	}
	if !strings.Contains(s.Email, "@") {
		return fmt.Errorf("%w: email must contain @", ErrInvalidStudent)
	}
	return nil
}

func ValidateCourse(c CourseGrade, seminar bool) error {
	if c.CourseCode == "" || c.CourseName == "" {
		return fmt.Errorf("%w: code and name are required", ErrInvalidCourse)
	}
	if c.Credits <= 0 || c.Credits > 12 {
		return fmt.Errorf("%w: credits out of range", ErrInvalidCourse)
	}
	if seminar {
		if c.Letter == "" && (math.IsNaN(c.Score) || c.Score < 0 || c.Score > 100) {
			return fmt.Errorf("%w: seminar requires letter or score", ErrInvalidCourse)
		}
		if c.Letter != "" && len(c.Letter) > 2 {
			return fmt.Errorf("%w: invalid letter", ErrInvalidCourse)
		}
	} else if math.IsNaN(c.Score) || c.Score < 0 || c.Score > 100 {
		return fmt.Errorf("%w: score out of range", ErrInvalidCourse)
	}
	return nil
}

func ValidateRecord(g Gradebook) error {
	if err := ValidateStudent(g.Student); err != nil {
		return err
	}
	if len(g.DegreeCourses) != 3 {
		return fmt.Errorf("%w: exactly three degree courses required", ErrInvalidCourse)
	}
	seen := map[string]bool{}
	for _, c := range g.DegreeCourses {
		if err := ValidateCourse(c, false); err != nil {
			return err
		}
		if seen[c.CourseCode] {
			return ErrDuplicateCourse
		}
		seen[c.CourseCode] = true
	}
	for _, c := range g.Seminars {
		if err := ValidateCourse(c, true); err != nil {
			return err
		}
		if seen[c.CourseCode] {
			return ErrDuplicateCourse
		}
		seen[c.CourseCode] = true
	}
	return nil
}
