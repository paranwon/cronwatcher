package watcher

import (
	"testing"
	"time"
)

func TestGetFlap_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newFlapStore()
	_, ok := s.GetFlap("missing")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestRecordFlapState_UnknownJob_ReturnsError(t *testing.T) {
	s := newFlapStore()
	err := s.RecordFlapState("missing", "success", time.Minute)
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestRecordFlapState_AlternatingStates_IncrementsCount(t *testing.T) {
	s := newFlapStore()
	s.registerFlapJob("backup")

	window := time.Minute
	_ = s.RecordFlapState("backup", "success", window)
	_ = s.RecordFlapState("backup", "failure", window)
	_ = s.RecordFlapState("backup", "success", window)

	entry, ok := s.GetFlap("backup")
	if !ok {
		t.Fatal("expected entry")
	}
	if entry.Flaps != 2 {
		t.Fatalf("expected 2 flaps, got %d", entry.Flaps)
	}
}

func TestRecordFlapState_SameState_DoesNotIncrement(t *testing.T) {
	s := newFlapStore()
	s.registerFlapJob("sync")

	window := time.Minute
	_ = s.RecordFlapState("sync", "success", window)
	_ = s.RecordFlapState("sync", "success", window)
	_ = s.RecordFlapState("sync", "success", window)

	entry, ok := s.GetFlap("sync")
	if !ok {
		t.Fatal("expected entry")
	}
	if entry.Flaps != 0 {
		t.Fatalf("expected 0 flaps, got %d", entry.Flaps)
	}
}

func TestRecordFlapState_WindowExpiry_ResetsCount(t *testing.T) {
	s := newFlapStore()
	s.registerFlapJob("cleanup")

	// Use a tiny window so it expires immediately.
	tinyWindow := time.Nanosecond

	_ = s.RecordFlapState("cleanup", "success", tinyWindow)
	_ = s.RecordFlapState("cleanup", "failure", tinyWindow)

	// Sleep to let window expire.
	time.Sleep(2 * time.Millisecond)

	_ = s.RecordFlapState("cleanup", "success", tinyWindow)

	entry, ok := s.GetFlap("cleanup")
	if !ok {
		t.Fatal("expected entry")
	}
	// After window reset the alternation from before should not count.
	if entry.Flaps != 0 {
		t.Fatalf("expected 0 flaps after window reset, got %d", entry.Flaps)
	}
}
