package watcher_test

import (
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestPauseJob_MarksJobPaused(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	if err := w.PauseJob("test-job"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !w.IsPaused("test-job") {
		t.Error("expected job to be paused")
	}
}

func TestResumeJob_ClearsPausedFlag(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	_ = w.PauseJob("test-job")
	if err := w.ResumeJob("test-job"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.IsPaused("test-job") {
		t.Error("expected job to not be paused after resume")
	}
}

func TestPauseJob_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	if err := w.PauseJob("no-such-job"); err == nil {
		t.Error("expected error for unknown job, got nil")
	}
}

func TestResumeJob_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	if err := w.ResumeJob("no-such-job"); err == nil {
		t.Error("expected error for unknown job, got nil")
	}
}

func TestIsPaused_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	if w.IsPaused("ghost-job") {
		t.Error("expected false for unknown job")
	}
}
