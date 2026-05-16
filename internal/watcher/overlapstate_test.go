package watcher_test

import (
	"testing"
	"time"

	"github.com/densestvoid/cronwatcher/internal/watcher"
)

func TestRecordOverlap_StoresEntry(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	runningSince := time.Now().Add(-30 * time.Second)
	if err := w.RecordOverlap("job1", runningSince); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, ok := w.GetOverlaps("job1")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !entries[0].RunSince.Equal(runningSince) {
		t.Errorf("expected RunSince %v, got %v", runningSince, entries[0].RunSince)
	}
	if entries[0].JobName != "job1" {
		t.Errorf("expected JobName job1, got %s", entries[0].JobName)
	}
}

func TestRecordOverlap_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	if err := w.RecordOverlap("nonexistent", time.Now()); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetOverlaps_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	_, ok := w.GetOverlaps("nonexistent")
	if ok {
		t.Fatal("expected ok=false for unknown job")
	}
}

func TestGetOverlaps_ReturnsCopy(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_ = w.RecordOverlap("job1", time.Now())

	entries, _ := w.GetOverlaps("job1")
	entries[0].JobName = "mutated"

	again, _ := w.GetOverlaps("job1")
	if again[0].JobName != "job1" {
		t.Error("GetOverlaps should return a copy, not a reference")
	}
}

func TestClearOverlaps_RemovesEntries(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_ = w.RecordOverlap("job1", time.Now())
	_ = w.RecordOverlap("job1", time.Now())

	if err := w.ClearOverlaps("job1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, ok := w.GetOverlaps("job1")
	if !ok {
		t.Fatal("expected ok=true after clear")
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(entries))
	}
}

func TestClearOverlaps_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	if err := w.ClearOverlaps("nonexistent"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}
