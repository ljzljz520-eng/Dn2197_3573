package ranking

import (
	"gradebook/fixtures"
	"gradebook/grading"
	"testing"
)

func TestRankingOrdersAndFilters(t *testing.T) {
	records := grading.SortByAverage(fixtures.Batch())
	if records[0].Student.Name != "Ada Lovelace" || records[0].Rank != 1 {
		t.Fatalf("unexpected ranking %#v", records)
	}
	if len(TopByRank(records, 2)) != 2 {
		t.Fatal("top limit failed")
	}
	if len(FindByCohort(records, 2025)) != 3 {
		t.Fatal("cohort filter failed")
	}
}
