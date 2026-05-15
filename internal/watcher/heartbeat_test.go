package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestRecordHeartbeat_StoresEntry(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	err := w.RecordHeartbeat("job1", 5*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, ok := w.GetHeartbeat("job1")
	if !ok {
		t.Fatal("expected heartbeat entry to exist")
	}
	if e.Missed {
		t.Error("expected Missed to be false immediately after heartbeat")
	}
	if e.ExpectsBy.Before(time.Now()) {
		t.Error("expected ExpectsBy to be in the future")
	}
}

func TestRecordHeartbeat_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	err := w.RecordHeartbeat("no-such-job", time.Minute)
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetHeartbeat_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_, ok := w.GetHeartbeat("missing")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestCheckHeartbeats_MarksMissedWhenExpired(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	// Record a heartbeat that expired immediately
	err := w.RecordHeartbeat("job1", -1*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	missed := w.CheckHeartbeats()
	if len(missed) != 1 || missed[0] != "job1" {
		t.Errorf("expected [job1] to be missed, got %v", missed)
	}

	e, _ := w.GetHeartbeat("job1")
	if !e.Missed {
		t.Error("expected Missed flag to be set")
	}
}

func TestCheckHeartbeats_DoesNotRepeatMissed(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_ = w.RecordHeartbeat("job1", -1*time.Second)

	first := w.CheckHeartbeats()
	if len(first) != 1 {
		t.Fatalf("expected 1 missed on first check, got %d", len(first))
	}

	second := w.CheckHeartbeats()
	if len(second) != 0 {
		t.Errorf("expected no missed on second check, got %v", second)
	}
}

func TestCheckHeartbeats_FutureDeadline_NotMissed(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_ = w.RecordHeartbeat("job1", 10*time.Minute)

	missed := w.CheckHeartbeats()
	if len(missed) != 0 {
		t.Errorf("expected no missed jobs, got %v", missed)
	}
}
