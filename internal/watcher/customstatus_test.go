package watcher

import (
	"testing"
)

func TestRecordCustomStatus_StoresEntry(t *testing.T) {
	w := NewForTest([]string{"backup"})
	if err := w.RecordCustomStatus("backup", "all good", "info"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := w.GetCustomStatus("backup")
	if !ok {
		t.Fatal("expected status to be present")
	}
	if s.Message != "all good" {
		t.Errorf("expected message 'all good', got %q", s.Message)
	}
	if s.Severity != "info" {
		t.Errorf("expected severity 'info', got %q", s.Severity)
	}
	if s.RecordedAt.IsZero() {
		t.Error("expected RecordedAt to be set")
	}
}

func TestRecordCustomStatus_UnknownJob_ReturnsError(t *testing.T) {
	w := NewForTest([]string{"backup"})
	if err := w.RecordCustomStatus("ghost", "msg", "warn"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetCustomStatus_UnknownJob_ReturnsFalse(t *testing.T) {
	w := NewForTest([]string{"backup"})
	_, ok := w.GetCustomStatus("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestClearCustomStatus_RemovesEntry(t *testing.T) {
	w := NewForTest([]string{"backup"})
	_ = w.RecordCustomStatus("backup", "degraded", "warn")
	if err := w.ClearCustomStatus("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := w.GetCustomStatus("backup")
	if ok {
		t.Fatal("expected status to be cleared")
	}
}

func TestClearCustomStatus_UnknownJob_ReturnsError(t *testing.T) {
	w := NewForTest([]string{"backup"})
	if err := w.ClearCustomStatus("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestRecordCustomStatus_OverwritesPrevious(t *testing.T) {
	w := NewForTest([]string{"backup"})
	_ = w.RecordCustomStatus("backup", "first", "info")
	_ = w.RecordCustomStatus("backup", "second", "error")
	s, ok := w.GetCustomStatus("backup")
	if !ok {
		t.Fatal("expected status")
	}
	if s.Message != "second" {
		t.Errorf("expected 'second', got %q", s.Message)
	}
	if s.Severity != "error" {
		t.Errorf("expected 'error', got %q", s.Severity)
	}
}
