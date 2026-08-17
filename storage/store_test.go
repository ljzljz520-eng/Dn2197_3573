package storage

import (
	"gradebook/domain"
	"path/filepath"
	"testing"
)

func TestStoreRoundTripAndFilters(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "grades.db"))
	if err != nil {
		t.Fatal(err)
	}
	record := domain.Gradebook{RecordID: "S1", Student: domain.NewStudent("S1", "Ada", "CS", 2025, "a@example.edu", "Dr"), DegreeCourses: []domain.CourseGrade{domain.NewCourseGrade("G1", "One", 3, 90, ""), domain.NewCourseGrade("G2", "Two", 3, 80, ""), domain.NewCourseGrade("G3", "Three", 3, 70, "")}}
	if err := store.Save(record); err != nil {
		t.Fatal(err)
	}
	if err := store.AppendAudit(domain.AuditEvent{ID: "a1", RecordID: "S1", Operator: "op", Action: "create", Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get("S1")
	if err != nil || got.Student.Name != "Ada" {
		t.Fatalf("round trip %v %#v", err, got)
	}
	if len(FilterActive([]domain.Gradebook{got})) != 1 || len(SortByName([]domain.Gradebook{got})) != 1 {
		t.Fatal("record helpers failed")
	}
	if err := store.SetMeta("schema", "1"); err != nil {
		t.Fatal(err)
	}
	if value, err := store.GetMeta("schema"); err != nil || value != "1" {
		t.Fatalf("meta %q %v", value, err)
	}
	if !IsSchemaReady(BucketNames()) {
		t.Fatal("schema not ready")
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
}
