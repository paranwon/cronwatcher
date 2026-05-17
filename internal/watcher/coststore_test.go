package watcher

import (
	"testing"
	"time"
)

func TestRecordCost_UnknownJob_ReturnsError(t *testing.T) {
	s := newCostStore([]string{"job-a"})
	err := s.recordCost("unknown", CostEntry{CPUSecs: 1.0})
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetCostSummary_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newCostStore([]string{"job-a"})
	_, ok := s.getCostSummary("unknown")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetCostSummary_NoEntries_ReturnsEmptySummary(t *testing.T) {
	s := newCostStore([]string{"job-a"})
	sum, ok := s.getCostSummary("job-a")
	if !ok {
		t.Fatal("expected true for known job")
	}
	if sum.RunCount != 0 {
		t.Fatalf("expected 0 runs, got %d", sum.RunCount)
	}
}

func TestRecordCost_StoresAndAverages(t *testing.T) {
	s := newCostStore([]string{"job-a"})
	now := time.Now()
	_ = s.recordCost("job-a", CostEntry{RecordedAt: now, CPUSecs: 2.0, MemoryMB: 100})
	_ = s.recordCost("job-a", CostEntry{RecordedAt: now, CPUSecs: 4.0, MemoryMB: 200})

	sum, ok := s.getCostSummary("job-a")
	if !ok {
		t.Fatal("expected summary")
	}
	if sum.RunCount != 2 {
		t.Fatalf("expected 2 runs, got %d", sum.RunCount)
	}
	if sum.TotalCPU != 6.0 {
		t.Fatalf("expected TotalCPU=6.0, got %f", sum.TotalCPU)
	}
	if sum.AvgCPU != 3.0 {
		t.Fatalf("expected AvgCPU=3.0, got %f", sum.AvgCPU)
	}
	if sum.TotalMemMB != 300 {
		t.Fatalf("expected TotalMemMB=300, got %f", sum.TotalMemMB)
	}
	if sum.AvgMemMB != 150 {
		t.Fatalf("expected AvgMemMB=150, got %f", sum.AvgMemMB)
	}
}
