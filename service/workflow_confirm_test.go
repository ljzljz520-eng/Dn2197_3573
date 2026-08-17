package service

import (
	"gradebook/fixtures"
	"gradebook/storage"
	"path/filepath"
	"sync"
	"testing"
)

func TestWorkflowConcurrentConfirmationKeepsBothOperators(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "grades.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	app := NewGradebookService(store)
	if err := app.Add(fixtures.Record("S1", "Ada Lovelace", "Computing", 96, 91, 89), "registrar"); err != nil {
		t.Fatal(err)
	}
	gate := NewConfirmationGate()
	var wg sync.WaitGroup
	errors := make(chan error, 2)
	for _, operator := range []string{"alice", "bob"} {
		wg.Add(1)
		go func(operator string) { defer wg.Done(); errors <- app.Confirm("S1", operator, gate) }(operator)
	}
	gate.AwaitBothReads()
	gate.ReleaseBoth()
	wg.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	count, err := app.ConfirmationCount("S1")
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("confirmation count = %d, want 2", count)
	}
	operators, err := app.ConfirmedBy("S1")
	if err != nil {
		t.Fatal(err)
	}
	if len(operators) != 2 {
		t.Errorf("confirmed operators = %#v, want two", operators)
	}
}
