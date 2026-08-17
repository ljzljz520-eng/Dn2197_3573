package service

import (
	"gradebook/fixtures"
	"gradebook/storage"
	"path/filepath"
	"testing"
)

func TestWorkflowAddAndQuery(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "grades.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	app := NewGradebookService(store)
	if err := app.Add(fixtures.Record("S1", "Ada Lovelace", "Computing", 96, 91, 89), "registrar"); err != nil {
		t.Fatal(err)
	}
	got, err := app.ByID("S1")
	if err != nil || got.Average <= 90 {
		t.Fatalf("query failed %#v %v", got, err)
	}
	byName, err := app.ByName("ada")
	if err != nil || len(byName) != 1 {
		t.Fatalf("name query failed %#v %v", byName, err)
	}
	events, err := store.Audits("S1")
	if err != nil || len(events) != 1 || events[0].Action != "create" {
		t.Fatalf("audit failed %#v %v", events, err)
	}
}
