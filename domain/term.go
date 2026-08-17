package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type TermPlan struct {
	StudentID      string
	Term           string
	Courses        []Course
	MaximumCredits int
	Approved       bool
	Approver       string
}

type PlanConflict struct {
	Code    string
	Message string
}

func NewTermPlan(studentID, term string, maximumCredits int) TermPlan {
	return TermPlan{
		StudentID:      strings.TrimSpace(studentID),
		Term:           strings.TrimSpace(term),
		Courses:        []Course{},
		MaximumCredits: maximumCredits,
	}
}

func (p *TermPlan) AddCourse(course Course) error {
	if strings.TrimSpace(course.Code) == "" {
		return errors.New("course code is required")
	}
	for _, existing := range p.Courses {
		if strings.EqualFold(existing.Code, course.Code) {
			return fmt.Errorf("course %s already planned", course.Code)
		}
	}
	if p.Credits()+course.Credits > p.MaximumCredits {
		return fmt.Errorf("course %s exceeds maximum credits", course.Code)
	}
	p.Courses = append(p.Courses, course)
	sort.Slice(p.Courses, func(i, j int) bool {
		return p.Courses[i].Code < p.Courses[j].Code
	})
	p.Approved = false
	p.Approver = ""
	return nil
}

func (p *TermPlan) RemoveCourse(code string) bool {
	needle := strings.ToUpper(strings.TrimSpace(code))
	for index, course := range p.Courses {
		if strings.ToUpper(course.Code) == needle {
			p.Courses = append(p.Courses[:index], p.Courses[index+1:]...)
			p.Approved = false
			p.Approver = ""
			return true
		}
	}
	return false
}

func (p TermPlan) Credits() int {
	total := 0
	for _, course := range p.Courses {
		total += course.Credits
	}
	return total
}

func (p TermPlan) Conflicts(catalog Catalog) []PlanConflict {
	conflicts := []PlanConflict{}
	if p.StudentID == "" {
		conflicts = append(conflicts, PlanConflict{Code: "student", Message: "student is required"})
	}
	if p.Term == "" {
		conflicts = append(conflicts, PlanConflict{Code: "term", Message: "term is required"})
	}
	if len(p.Courses) == 0 {
		conflicts = append(conflicts, PlanConflict{Code: "courses", Message: "at least one course is required"})
	}
	if p.Credits() > p.MaximumCredits {
		conflicts = append(conflicts, PlanConflict{Code: "credits", Message: "planned credits exceed the maximum"})
	}
	for _, course := range p.Courses {
		catalogCourse, found := catalog.Lookup(course.Code)
		if !found {
			conflicts = append(conflicts, PlanConflict{Code: "catalog", Message: fmt.Sprintf("course %s is not in the catalog", course.Code)})
			continue
		}
		if catalogCourse.Credits != course.Credits {
			conflicts = append(conflicts, PlanConflict{Code: "credits", Message: fmt.Sprintf("course %s has mismatched credits", course.Code)})
		}
	}
	sort.SliceStable(conflicts, func(i, j int) bool {
		if conflicts[i].Code == conflicts[j].Code {
			return conflicts[i].Message < conflicts[j].Message
		}
		return conflicts[i].Code < conflicts[j].Code
	})
	return conflicts
}

func (p *TermPlan) Approve(approver string, catalog Catalog) error {
	name := strings.TrimSpace(approver)
	if name == "" {
		return errors.New("approver is required")
	}
	conflicts := p.Conflicts(catalog)
	if len(conflicts) > 0 {
		return fmt.Errorf("term plan has %d conflict(s)", len(conflicts))
	}
	p.Approved = true
	p.Approver = name
	return nil
}

func PlanCourseCodes(plan TermPlan) []string {
	codes := make([]string, 0, len(plan.Courses))
	for _, course := range plan.Courses {
		codes = append(codes, course.Code)
	}
	sort.Strings(codes)
	return codes
}

func PlanSummary(plan TermPlan) string {
	status := "pending"
	if plan.Approved {
		status = "approved by " + plan.Approver
	}
	return fmt.Sprintf("%s %s: %d courses, %d credits, %s", plan.StudentID, plan.Term, len(plan.Courses), plan.Credits(), status)
}
