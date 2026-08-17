package service

import (
	"gradebook/fixtures"
	"gradebook/storage"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "grades.db")
	store, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	app := NewGradebookService(store)
	if err := app.Add(fixtures.WithSeminar(fixtures.Record("S1", "Ada Lovelace", "Computing", 96, 91, 89), fixtures.Seminar("SEM1", "Research Seminar", "A", 0)), "registrar"); err != nil {
		t.Fatal(err)
	}
	if err := app.Confirm("S1", "alice", nil); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	reopenedApp := NewGradebookService(reopened)
	record, err := reopenedApp.Get("S1")
	if err != nil || record.ConfirmationCount != 1 || len(record.Seminars) != 1 {
		t.Fatalf("reopen failed %#v %v", record, err)
	}
	events, err := reopened.Audits("S1")
	if err != nil || len(events) != 2 {
		t.Fatalf("audit reopen failed %#v %v", events, err)
	}
}
