package watcher

import (
	"testing"
	"time"
)

func TestGetWindow_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newWindowStore([]string{"job-a"})
	_, ok := s.GetWindow("missing", time.Hour)
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestRecordWindow_UnknownJob_ReturnsError(t *testing.T) {
	s := newWindowStore([]string{"job-a"})
	err := s.RecordWindow("missing", true, time.Second)
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetWindow_EmptyReturnsZeroCounts(t *testing.T) {
	s := newWindowStore([]string{"job-a"})
	e, ok := s.GetWindow("job-a", time.Hour)
	if !ok {
		t.Fatal("expected ok for known job")
	}
	if e.Runs != 0 || e.Successes != 0 || e.Failures != 0 {
		t.Fatalf("expected zero counts, got %+v", e)
	}
}

func TestGetWindow_CountsSuccessAndFailure(t *testing.T) {
	s := newWindowStore([]string{"job-a"})
	_ = s.RecordWindow("job-a", true, 200*time.Millisecond)
	_ = s.RecordWindow("job-a", true, 400*time.Millisecond)
	_ = s.RecordWindow("job-a", false, 100*time.Millisecond)

	e, ok := s.GetWindow("job-a", time.Hour)
	if !ok {
		t.Fatal("expected ok")
	}
	if e.Runs != 3 {
		t.Errorf("expected 3 runs, got %d", e.Runs)
	}
	if e.Successes != 2 {
		t.Errorf("expected 2 successes, got %d", e.Successes)
	}
	if e.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", e.Failures)
	}
	expectedAvg := (200.0 + 400.0 + 100.0) / 3.0
	if e.AvgMs != expectedAvg {
		t.Errorf("expected avg %.2f, got %.2f", expectedAvg, e.AvgMs)
	}
}

func TestGetWindow_ExcludesOldEntries(t *testing.T) {
	s := newWindowStore([]string{"job-a"})
	// Manually inject an old record
	s.mu.Lock()
	s.records["job-a"] = append(s.records["job-a"], windowRecord{
		ts:      time.Now().Add(-2 * time.Hour),
		success: true,
		durMs:   500,
	})
	s.mu.Unlock()

	_ = s.RecordWindow("job-a", true, 100*time.Millisecond)

	e, ok := s.GetWindow("job-a", time.Hour)
	if !ok {
		t.Fatal("expected ok")
	}
	if e.Runs != 1 {
		t.Errorf("expected 1 run within window, got %d", e.Runs)
	}
}
