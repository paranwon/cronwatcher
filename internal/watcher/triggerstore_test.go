package watcher

import (
	"testing"
)

func TestRecordTrigger_StoresEntry(t *testing.T) {
	s := newTriggerStore()
	s.register("backup")

	err := s.recordTrigger("backup", "admin", "manual run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, ok := s.getTrigger("backup")
	if !ok {
		t.Fatal("expected trigger entry, got none")
	}
	if e.TriggeredBy != "admin" {
		t.Errorf("expected triggered_by=admin, got %s", e.TriggeredBy)
	}
	if e.Reason != "manual run" {
		t.Errorf("expected reason='manual run', got %s", e.Reason)
	}
	if e.TriggeredAt.IsZero() {
		t.Error("expected non-zero triggered_at")
	}
}

func TestRecordTrigger_UnknownJob_ReturnsError(t *testing.T) {
	s := newTriggerStore()
	err := s.recordTrigger("ghost", "admin", "")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetTrigger_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newTriggerStore()
	_, ok := s.getTrigger("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetTrigger_NoEntryYet_ReturnsFalse(t *testing.T) {
	s := newTriggerStore()
	s.register("backup")
	_, ok := s.getTrigger("backup")
	if ok {
		t.Fatal("expected false before any trigger recorded")
	}
}

func TestClearTrigger_RemovesEntry(t *testing.T) {
	s := newTriggerStore()
	s.register("backup")
	_ = s.recordTrigger("backup", "ci", "scheduled")

	err := s.clearTrigger("backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.getTrigger("backup")
	if ok {
		t.Fatal("expected no trigger after clear")
	}
}

func TestClearTrigger_UnknownJob_ReturnsError(t *testing.T) {
	s := newTriggerStore()
	err := s.clearTrigger("ghost")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}
