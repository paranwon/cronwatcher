package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestSetNextRun_StoresEntry(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	next := time.Now().Add(5 * time.Minute)

	err := w.SetNextRun("backup", next, "0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, ok := w.GetNextRun("backup")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if !e.NextRun.Equal(next) {
		t.Errorf("expected NextRun %v, got %v", next, e.NextRun)
	}
	if e.Schedule != "0 * * * *" {
		t.Errorf("expected schedule '0 * * * *', got %q", e.Schedule)
	}
	if e.JobName != "backup" {
		t.Errorf("expected job name 'backup', got %q", e.JobName)
	}
}

func TestSetNextRun_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	err := w.SetNextRun("ghost", time.Now(), "* * * * *")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetNextRun_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_, ok := w.GetNextRun("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestOverdueJobs_ReturnsExpiredEntries(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	past := time.Now().Add(-10 * time.Minute)
	future := time.Now().Add(10 * time.Minute)

	_ = w.SetNextRun("backup", past, "0 * * * *")
	_ = w.SetNextRun("cleanup", future, "0 0 * * *")

	overdue := w.OverdueJobs(time.Now())
	if len(overdue) != 1 {
		t.Fatalf("expected 1 overdue job, got %d", len(overdue))
	}
	if overdue[0].JobName != "backup" {
		t.Errorf("expected 'backup' to be overdue, got %q", overdue[0].JobName)
	}
}

func TestOverdueJobs_NoneOverdue_ReturnsEmpty(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_ = w.SetNextRun("backup", time.Now().Add(1*time.Hour), "0 * * * *")

	overdue := w.OverdueJobs(time.Now())
	if len(overdue) != 0 {
		t.Errorf("expected no overdue jobs, got %d", len(overdue))
	}
}
