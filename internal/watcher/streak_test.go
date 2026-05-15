package watcher_test

import (
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestGetStreak_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_, ok := w.GetStreak("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestRecordStreakSuccess_IncrementsSuccess(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	if err := w.RecordStreakSuccess("job1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.RecordStreakSuccess("job1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := w.GetStreak("job1")
	if !ok {
		t.Fatal("expected streak entry")
	}
	if e.SuccessStreak != 2 {
		t.Errorf("expected SuccessStreak=2, got %d", e.SuccessStreak)
	}
	if e.FailureStreak != 0 {
		t.Errorf("expected FailureStreak=0, got %d", e.FailureStreak)
	}
}

func TestRecordStreakFailure_IncrementsFailure(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_ = w.RecordStreakSuccess("job1")
	_ = w.RecordStreakFailure("job1")
	e, _ := w.GetStreak("job1")
	if e.FailureStreak != 1 {
		t.Errorf("expected FailureStreak=1, got %d", e.FailureStreak)
	}
	if e.SuccessStreak != 0 {
		t.Errorf("expected SuccessStreak reset to 0, got %d", e.SuccessStreak)
	}
}

func TestRecordStreakSuccess_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	if err := w.RecordStreakSuccess("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestRecordStreakFailure_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	if err := w.RecordStreakFailure("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}
