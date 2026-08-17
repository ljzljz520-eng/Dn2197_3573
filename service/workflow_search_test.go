package service

import (
	"gradebook/fixtures"
	"gradebook/storage"
	"path/filepath"
	"testing"
)

func TestWorkflowRankingStatisticsAndReport(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "grades.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	app := NewGradebookService(store)
	for _, record := range fixtures.Batch() {
		if err := app.Add(record, "registrar"); err != nil {
			t.Fatal(err)
		}
	}
	ranked, err := app.Ranked(2)
	if err != nil || len(ranked) != 2 || ranked[0].Rank != 1 {
		t.Fatalf("rank failed %#v %v", ranked, err)
	}
	stats, err := app.Statistics()
	if err != nil || stats.Count != 3 || stats.Highest < stats.Lowest {
		t.Fatalf("stats failed %#v %v", stats, err)
	}
	report, err := app.BuildReport("Term report")
	if err != nil || len(report.Lines) != 3 {
		t.Fatalf("report failed %#v %v", report, err)
	}
	if len(RenderReport(report)) == 0 {
		t.Fatal("empty report")
	}
}
