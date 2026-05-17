package watcher

import (
	"testing"
)

func TestRecordBlocked_StoresEntry(t *testing.T) {
	s := newBlockedStore([]string{"backup"})
	err := s.recordBlocked("backup", "migration", "lock held")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, ok := s.getBlocked("backup")
	if !ok {
		t.Fatal("expected entries to exist")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].BlockedBy != "migration" {
		t.Errorf("expected BlockedBy=migration, got %s", entries[0].BlockedBy)
	}
	if entries[0].Reason != "lock held" {
		t.Errorf("expected Reason='lock held', got %s", entries[0].Reason)
	}
}

func TestRecordBlocked_UnknownJob_ReturnsError(t *testing.T) {
	s := newBlockedStore([]string{})
	err := s.recordBlocked("nonexistent", "other", "reason")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetBlocked_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newBlockedStore([]string{})
	_, ok := s.getBlocked("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetBlocked_ReturnsCopy(t *testing.T) {
	s := newBlockedStore([]string{"job1"})
	_ = s.recordBlocked("job1", "job2", "waiting")
	entries, _ := s.getBlocked("job1")
	entries[0].Reason = "mutated"
	orig, _ := s.getBlocked("job1")
	if orig[0].Reason == "mutated" {
		t.Error("getBlocked should return a copy, not a reference")
	}
}

func TestClearBlocked_RemovesEntries(t *testing.T) {
	s := newBlockedStore([]string{"job1"})
	_ = s.recordBlocked("job1", "job2", "reason")
	if err := s.clearBlocked("job1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, ok := s.getBlocked("job1")
	if !ok {
		t.Fatal("job1 should still be known")
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(entries))
	}
}

func TestClearBlocked_UnknownJob_ReturnsError(t *testing.T) {
	s := newBlockedStore([]string{})
	err := s.clearBlocked("ghost")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}
