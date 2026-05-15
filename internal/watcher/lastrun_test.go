package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestGetLastRun_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_, ok := w.GetLastRun("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestRecordLastRun_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	err := w.RecordLastRun("ghost", watcher.LastRunEntry{})
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestRecordLastRun_StoresEntry(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	now := time.Now()
	entry := watcher.LastRunEntry{
		FinishedAt: now,
		Duration:   5 * time.Second,
		Success:    true,
	}
	if err := w.RecordLastRun("test-job", entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := w.GetLastRun("test-job")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if !got.FinishedAt.Equal(now) {
		t.Errorf("expected FinishedAt %v, got %v", now, got.FinishedAt)
	}
	if got.Duration != 5*time.Second {
		t.Errorf("expected duration 5s, got %v", got.Duration)
	}
	if !got.Success {
		t.Error("expected Success to be true")
	}
}

func TestRecordLastRun_OverwritesPrevious(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	first := watcher.LastRunEntry{FinishedAt: time.Now().Add(-1 * time.Hour), Success: false, Error: "timeout"}
	second := watcher.LastRunEntry{FinishedAt: time.Now(), Success: true}

	_ = w.RecordLastRun("test-job", first)
	_ = w.RecordLastRun("test-job", second)

	got, ok := w.GetLastRun("test-job")
	if !ok {
		t.Fatal("expected entry")
	}
	if !got.Success {
		t.Error("expected second entry to overwrite first")
	}
	if got.Error != "" {
		t.Errorf("expected empty error, got %q", got.Error)
	}
}
