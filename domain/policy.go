package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type GradePolicy struct {
	MinimumDegreeScore    float64
	MinimumAverage        float64
	RequiredDegreeCourses int
	RequiredSeminars      int
	MinimumCredits        int
	AllowedLetters        []string
}

type PolicyViolation struct {
	Code       string
	Message    string
	CourseCode string
}

func DefaultGradePolicy() GradePolicy {
	return GradePolicy{
		MinimumDegreeScore:    60,
		MinimumAverage:        60,
		RequiredDegreeCourses: 3,
		RequiredSeminars:      1,
		MinimumCredits:        10,
		AllowedLetters:        []string{"A+", "A", "A-", "B+", "B", "B-", "C+", "C"},
	}
}

func ValidatePolicy(policy GradePolicy) error {
	if policy.MinimumDegreeScore < 0 || policy.MinimumDegreeScore > 100 {
		return errors.New("minimum degree score out of range")
	}
	if policy.MinimumAverage < 0 || policy.MinimumAverage > 100 {
		return errors.New("minimum average out of range")
	}
	if policy.RequiredDegreeCourses < 0 {
		return errors.New("required degree courses cannot be negative")
	}
	if policy.RequiredSeminars < 0 {
		return errors.New("required seminars cannot be negative")
	}
	if policy.MinimumCredits < 0 {
		return errors.New("minimum credits cannot be negative")
	}
	if len(policy.AllowedLetters) == 0 {
		return errors.New("allowed letters cannot be empty")
	}
	return nil
}

func EvaluatePolicy(record Gradebook, policy GradePolicy) []PolicyViolation {
	violations := []PolicyViolation{}
	if err := ValidatePolicy(policy); err != nil {
		return append(violations, PolicyViolation{Code: "invalid_policy", Message: err.Error()})
	}
	if len(record.DegreeCourses) < policy.RequiredDegreeCourses {
		violations = append(violations, PolicyViolation{Code: "degree_count", Message: fmt.Sprintf("requires %d degree courses", policy.RequiredDegreeCourses)})
	}
	if len(record.Seminars) < policy.RequiredSeminars {
		violations = append(violations, PolicyViolation{Code: "seminar_count", Message: fmt.Sprintf("requires %d seminars", policy.RequiredSeminars)})
	}
	credits := 0
	for _, course := range record.DegreeCourses {
		credits += course.Credits
		if course.Score < policy.MinimumDegreeScore {
			violations = append(violations, PolicyViolation{Code: "degree_score", Message: "degree score below minimum", CourseCode: course.CourseCode})
		}
	}
	for _, course := range record.Seminars {
		credits += course.Credits
		if course.Letter != "" && !containsLetter(policy.AllowedLetters, course.Letter) {
			violations = append(violations, PolicyViolation{Code: "seminar_letter", Message: "seminar letter is not allowed", CourseCode: course.CourseCode})
		}
	}
	if credits < policy.MinimumCredits {
		violations = append(violations, PolicyViolation{Code: "credits", Message: fmt.Sprintf("requires %d credits", policy.MinimumCredits)})
	}
	if record.Average < policy.MinimumAverage {
		violations = append(violations, PolicyViolation{Code: "average", Message: "average below minimum"})
	}
	sort.SliceStable(violations, func(i, j int) bool {
		if violations[i].Code == violations[j].Code {
			return violations[i].CourseCode < violations[j].CourseCode
		}
		return violations[i].Code < violations[j].Code
	})
	return violations
}

func containsLetter(allowed []string, value string) bool {
	for _, letter := range allowed {
		if strings.EqualFold(letter, value) {
			return true
		}
	}
	return false
}

func PolicyPassed(record Gradebook, policy GradePolicy) bool {
	return len(EvaluatePolicy(record, policy)) == 0
}

func ViolationMessages(violations []PolicyViolation) []string {
	result := make([]string, 0, len(violations))
	for _, violation := range violations {
		if violation.CourseCode == "" {
			result = append(result, violation.Message)
		} else {
			result = append(result, violation.CourseCode+": "+violation.Message)
		}
	}
	return result
}

func MergePolicies(base, override GradePolicy) GradePolicy {
	result := base
	if override.MinimumDegreeScore != 0 {
		result.MinimumDegreeScore = override.MinimumDegreeScore
	}
	if override.MinimumAverage != 0 {
		result.MinimumAverage = override.MinimumAverage
	}
	if override.RequiredDegreeCourses != 0 {
		result.RequiredDegreeCourses = override.RequiredDegreeCourses
	}
	if override.RequiredSeminars != 0 {
		result.RequiredSeminars = override.RequiredSeminars
	}
	if override.MinimumCredits != 0 {
		result.MinimumCredits = override.MinimumCredits
	}
	if len(override.AllowedLetters) != 0 {
		result.AllowedLetters = append([]string(nil), override.AllowedLetters...)
	}
	return result
}
