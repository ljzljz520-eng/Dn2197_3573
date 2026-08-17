package grading

import (
	"gradebook/fixtures"
	"testing"
)

func TestWeightedAverageAndBands(t *testing.T) {
	record := fixtures.WithSeminar(fixtures.Record("S1", "Ada", "CS", 100, 90, 80), fixtures.Seminar("SEM", "Seminar", "A", 0))
	if got := WeightedAverage(record); got != 92 {
		t.Fatalf("unexpected average %.2f", got)
	}
	if PercentageBand(95) != "distinction" || PercentageBand(75) != "pass" || PercentageBand(40) != "review" {
		t.Fatal("unexpected bands")
	}
}
