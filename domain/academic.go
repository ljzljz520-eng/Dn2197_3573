package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Enrollment struct {
	StudentID  string
	CourseCode string
	Term       string
	Status     string
}

type Milestone struct {
	Name      string
	Completed bool
	Evidence  string
}

type AcademicProfile struct {
	Student     Student
	Enrollments []Enrollment
	Milestones  []Milestone
	Notes       []string
}

func NewEnrollment(studentID, courseCode, term string) Enrollment {
	return Enrollment{StudentID: strings.TrimSpace(studentID), CourseCode: strings.ToUpper(strings.TrimSpace(courseCode)), Term: strings.TrimSpace(term), Status: "active"}
}

func (e Enrollment) IsActive() bool { return strings.EqualFold(e.Status, "active") }

func CompleteMilestone(name, evidence string) Milestone {
	return Milestone{Name: strings.TrimSpace(name), Completed: true, Evidence: strings.TrimSpace(evidence)}
}

func ValidateEnrollment(e Enrollment) error {
	if e.StudentID == "" {
		return errors.New("enrollment student is required")
	}
	if e.CourseCode == "" {
		return errors.New("enrollment course is required")
	}
	if e.Term == "" {
		return errors.New("enrollment term is required")
	}
	if e.Status != "active" && e.Status != "completed" && e.Status != "withdrawn" {
		return fmt.Errorf("invalid enrollment status %q", e.Status)
	}
	return nil
}

func ValidateProfile(profile AcademicProfile) error {
	if err := ValidateStudent(profile.Student); err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, enrollment := range profile.Enrollments {
		if err := ValidateEnrollment(enrollment); err != nil {
			return err
		}
		key := enrollment.CourseCode + "@" + enrollment.Term
		if seen[key] {
			return fmt.Errorf("duplicate enrollment %s", key)
		}
		seen[key] = true
	}
	return nil
}

func ActiveEnrollments(profile AcademicProfile) []Enrollment {
	result := []Enrollment{}
	for _, enrollment := range profile.Enrollments {
		if enrollment.IsActive() {
			result = append(result, enrollment)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CourseCode < result[j].CourseCode })
	return result
}

func CompletedMilestones(profile AcademicProfile) []Milestone {
	result := []Milestone{}
	for _, milestone := range profile.Milestones {
		if milestone.Completed {
			result = append(result, milestone)
		}
	}
	return result
}

func AddNote(profile *AcademicProfile, note string) {
	value := strings.TrimSpace(note)
	if value != "" {
		profile.Notes = append(profile.Notes, value)
	}
}

func Progress(profile AcademicProfile) float64 {
	if len(profile.Milestones) == 0 {
		return 0
	}
	completed := len(CompletedMilestones(profile))
	return float64(completed) / float64(len(profile.Milestones)) * 100
}

func ProfileSummary(profile AcademicProfile) string {
	return fmt.Sprintf("%s: %d active enrollments, %.0f%% milestone progress", profile.Student.Name, len(ActiveEnrollments(profile)), Progress(profile))
}
