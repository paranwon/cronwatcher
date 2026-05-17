package watcher

import (
	"testing"
)

func TestRecordRetry_StoresEntry(t *testing.T) {
	s := newRetryStore()
	s.entries["backup"] = nil

	if err := s.recordRetry("backup", "timeout"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.entries["backup"]) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s.entries["backup"]))
	}
	if s.entries["backup"][0].Reason != "timeout" {
		t.Errorf("expected reason 'timeout', got %q", s.entries["backup"][0].Reason)
	}
	if s.entries["backup"][0].Attempt != 1 {
		t.Errorf("expected attempt 1, got %d", s.entries["backup"][0].Attempt)
	}
}

func TestRecordRetry_UnknownJob_ReturnsError(t *testing.T) {
	s := newRetryStore()

	err := s.recordRetry("ghost", "fail")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetRetrySummary_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newRetryStore()

	_, ok := s.getRetrySummary("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetRetrySummary_ReturnsCorrectCounts(t *testing.T) {
	s := newRetryStore()
	s.entries["sync"] = nil

	_ = s.recordRetry("sync", "network error")
	_ = s.recordRetry("sync", "timeout")

	summary, ok := s.getRetrySummary("sync")
	if !ok {
		t.Fatal("expected summary to exist")
	}
	if summary.Attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", summary.Attempts)
	}
	if summary.LastReason != "timeout" {
		t.Errorf("expected last reason 'timeout', got %q", summary.LastReason)
	}
}

func TestClearRetries_RemovesEntries(t *testing.T) {
	s := newRetryStore()
	s.entries["etl"] = nil
	_ = s.recordRetry("etl", "crash")

	if err := s.clearRetries("etl"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	summary, ok := s.getRetrySummary("etl")
	if !ok {
		t.Fatal("expected summary to still exist after clear")
	}
	if summary.Attempts != 0 {
		t.Errorf("expected 0 attempts after clear, got %d", summary.Attempts)
	}
}

func TestClearRetries_UnknownJob_ReturnsError(t *testing.T) {
	s := newRetryStore()

	if err := s.clearRetries("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}
