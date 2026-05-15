package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatcher/internal/watcher"
)

func TestGetAvgDuration_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest([]string{"job-a"})
	_, ok := w.GetAvgDuration("nonexistent")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestRecordDuration_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest([]string{"job-a"})
	err := w.RecordDuration("nonexistent", time.Second)
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestRecordDuration_StoresAndAverages(t *testing.T) {
	w := watcher.NewForTest([]string{"job-a"})

	if err := w.RecordDuration("job-a", 10*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.RecordDuration("job-a", 20*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, ok := w.GetAvgDuration("job-a")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Count != 2 {
		t.Errorf("expected count 2, got %d", e.Count)
	}
	if e.Total != 30*time.Second {
		t.Errorf("expected total 30s, got %v", e.Total)
	}
	if e.Average != 15*time.Second {
		t.Errorf("expected average 15s, got %v", e.Average)
	}
}

func TestRecordDuration_SingleEntry_AverageEqualsValue(t *testing.T) {
	w := watcher.NewForTest([]string{"job-b"})

	if err := w.RecordDuration("job-b", 5*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, ok := w.GetAvgDuration("job-b")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Average != 5*time.Second {
		t.Errorf("expected average 5s, got %v", e.Average)
	}
}
