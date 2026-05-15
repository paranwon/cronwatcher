package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestRecordTimeout_StoresEntry(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	if err := w.RecordTimeout("backup", 5*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecordTimeout_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	err := w.RecordTimeout("ghost", time.Minute)
	if err == nil {
		t.Fatal("expected error for unknown job, got nil")
	}
}

func TestCheckTimeouts_ReturnsOverdueJobs(t *testing.T) {
	cfg := testConfig()
	w := watcher.NewForTest(cfg)

	// Record a timeout that started 10 minutes ago with a 5 minute max.
	w.AdvanceClock(-10 * time.Minute)
	if err := w.RecordTimeout("backup", 5*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.AdvanceClock(10 * time.Minute) // restore to "now"

	overdue := w.CheckTimeouts()
	if len(overdue) != 1 || overdue[0] != "backup" {
		t.Fatalf("expected [backup], got %v", overdue)
	}
}

func TestCheckTimeouts_WithinLimit_NotReturned(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	if err := w.RecordTimeout("backup", 10*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	overdue := w.CheckTimeouts()
	if len(overdue) != 0 {
		t.Fatalf("expected no overdue jobs, got %v", overdue)
	}
}

func TestClearTimeout_RemovesEntry(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	w.AdvanceClock(-10 * time.Minute)
	_ = w.RecordTimeout("backup", 5*time.Minute)
	w.AdvanceClock(10 * time.Minute)

	w.ClearTimeout("backup")

	overdue := w.CheckTimeouts()
	if len(overdue) != 0 {
		t.Fatalf("expected no overdue jobs after clear, got %v", overdue)
	}
}
