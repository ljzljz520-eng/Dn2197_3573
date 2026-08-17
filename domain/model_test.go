package domain

import "testing"

func TestCourseNormalization(t *testing.T) {
	record := Gradebook{Student: NewStudent(" S1 ", " Ada ", " CS ", 2025, "ADA@EXAMPLE.EDU", " Dr "), DegreeCourses: []CourseGrade{NewCourseGrade("g1", " One ", 3, 90, ""), NewCourseGrade("g2", "Two", 3, 80, ""), NewCourseGrade("g3", "Three", 3, 70, "")}}
	Normalize(&record)
	if record.Student.ID != "S1" || record.Student.Email != "ada@example.edu" || record.DegreeCourses[0].CourseCode != "G1" {
		t.Fatalf("normalization failed: %#v", record)
	}
}

func TestValidateRecordRejectsDuplicateCourse(t *testing.T) {
	record := Gradebook{Student: NewStudent("S1", "Ada", "CS", 2025, "a@example.edu", "Dr"), DegreeCourses: []CourseGrade{NewCourseGrade("G1", "One", 3, 90, ""), NewCourseGrade("G1", "Two", 3, 80, ""), NewCourseGrade("G3", "Three", 3, 70, "")}}
	if err := ValidateRecord(record); err != ErrDuplicateCourse {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}
