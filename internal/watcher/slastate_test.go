package watcher

import (
	"testing"
	"time"
)

func TestRecordSLA_StoresEntry(t *testing.T) {
	s := newSLAStore()
	deadline := time.Now().Add(time.Hour)
	if err := s.RecordSLA("backup", deadline, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.GetSLA("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if !e.MetSLA {
		t.Error("expected MetSLA to be true")
	}
	if !e.Deadline.Equal(deadline) {
		t.Errorf("expected deadline %v, got %v", deadline, e.Deadline)
	}
}

func TestGetSLA_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newSLAStore()
	_, ok := s.GetSLA("nonexistent")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestCheckSLAViolations_ReturnsViolations(t *testing.T) {
	s := newSLAStore()
	_ = s.RecordSLA("job-a", time.Now(), false)
	_ = s.RecordSLA("job-b", time.Now(), true)
	_ = s.RecordSLA("job-c", time.Now(), false)

	violations := s.CheckSLAViolations()
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestCheckSLAViolations_NoViolations_ReturnsEmpty(t *testing.T) {
	s := newSLAStore()
	_ = s.RecordSLA("job-a", time.Now(), true)
	_ = s.RecordSLA("job-b", time.Now(), true)

	violations := s.CheckSLAViolations()
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestRecordSLA_OverwritesPrevious(t *testing.T) {
	s := newSLAStore()
	deadline := time.Now().Add(time.Hour)
	_ = s.RecordSLA("backup", deadline, false)
	_ = s.RecordSLA("backup", deadline, true)

	e, ok := s.GetSLA("backup")
	if !ok {
		t.Fatal("expected entry")
	}
	if !e.MetSLA {
		t.Error("expected updated MetSLA to be true")
	}
}
