package watcher

import (
	"testing"
	"time"
)

func TestMarkStale_ReturnsStaleJobs(t *testing.T) {
	w := NewForTest(testConfig())
	now := time.Now()

	// Simulate a finish that happened a long time ago.
	oldFinish := now.Add(-10 * time.Minute)
	w.RecordStart("backup")
	w.RecordFinish("backup", nil)

	// Manually set LastSeen to simulate staleness.
	w.mu.Lock()
	s := w.jobs["backup"]
	s.LastSeen = oldFinish
	w.jobs["backup"] = s
	w.mu.Unlock()

	stale, err := w.MarkStale(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stale) != 1 {
		t.Fatalf("expected 1 stale job, got %d", len(stale))
	}
	if stale[0].JobName != "backup" {
		t.Errorf("expected stale job 'backup', got %q", stale[0].JobName)
	}
}

func TestMarkStale_SkipsRunningJobs(t *testing.T) {
	w := NewForTest(testConfig())
	now := time.Now()

	w.RecordStart("backup")
	// Job is still running — should not be considered stale.

	stale, err := w.MarkStale(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stale) != 0 {
		t.Errorf("expected no stale jobs, got %d", len(stale))
	}
}

func TestMarkStale_SkipsPausedJobs(t *testing.T) {
	w := NewForTest(testConfig())
	now := time.Now()

	w.RecordStart("backup")
	w.RecordFinish("backup", nil)
	w.PauseJob("backup")

	w.mu.Lock()
	s := w.jobs["backup"]
	s.LastSeen = now.Add(-10 * time.Minute)
	w.jobs["backup"] = s
	w.mu.Unlock()

	stale, err := w.MarkStale(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stale) != 0 {
		t.Errorf("expected no stale jobs for paused job, got %d", len(stale))
	}
}

func TestMarkStale_NoLastSeen_Skipped(t *testing.T) {
	w := NewForTest(testConfig())
	now := time.Now()

	// Never recorded a finish — LastSeen is zero.
	stale, err := w.MarkStale(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stale) != 0 {
		t.Errorf("expected no stale jobs when LastSeen is zero, got %d", len(stale))
	}
}
