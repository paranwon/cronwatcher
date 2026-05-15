package watcher_test

import (
	"testing"
	"time"

	"github.com/dop251/cronwatcher/internal/watcher"
)

func TestPruneHistory_RemovesOldEntries(t *testing.T) {
	cfg := testConfig()
	w := watcher.NewForTest(cfg)

	now := time.Now()

	// Simulate two runs: one old, one recent
	_ = w.RecordStart("backup", now.Add(-3*time.Hour))
	_ = w.RecordFinish("backup", now.Add(-3*time.Hour).Add(10*time.Second), nil)

	_ = w.RecordStart("backup", now.Add(-30*time.Minute))
	_ = w.RecordFinish("backup", now.Add(-30*time.Minute).Add(5*time.Second), nil)

	removed := w.PruneHistory(2 * time.Hour)
	if removed != 1 {
		t.Fatalf("expected 1 entry pruned, got %d", removed)
	}

	entries, ok := w.GetHistory("backup", 10)
	if !ok {
		t.Fatal("expected history to exist")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", len(entries))
	}
}

func TestPruneHistory_NothingToRemove(t *testing.T) {
	cfg := testConfig()
	w := watcher.NewForTest(cfg)

	now := time.Now()
	_ = w.RecordStart("backup", now.Add(-10*time.Minute))
	_ = w.RecordFinish("backup", now.Add(-10*time.Minute).Add(5*time.Second), nil)

	removed := w.PruneHistory(2 * time.Hour)
	if removed != 0 {
		t.Fatalf("expected 0 entries pruned, got %d", removed)
	}
}

func TestPruneHistory_EmptyWatcher_DoesNotPanic(t *testing.T) {
	cfg := testConfig()
	w := watcher.NewForTest(cfg)
	removed := w.PruneHistory(1 * time.Hour)
	if removed != 0 {
		t.Fatalf("expected 0, got %d", removed)
	}
}
