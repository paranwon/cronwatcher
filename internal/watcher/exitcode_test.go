package watcher

import (
	"testing"
)

func TestRecordExitCode_StoresEntry(t *testing.T) {
	w := NewForTest(testConfig())
	if err := w.RecordExitCode("backup", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := w.GetExitCode("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Code != 0 {
		t.Errorf("expected code 0, got %d", e.Code)
	}
	if e.RecordedAt.IsZero() {
		t.Error("expected RecordedAt to be set")
	}
}

func TestRecordExitCode_NonZeroCode(t *testing.T) {
	w := NewForTest(testConfig())
	if err := w.RecordExitCode("backup", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := w.GetExitCode("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Code != 1 {
		t.Errorf("expected code 1, got %d", e.Code)
	}
}

func TestRecordExitCode_UnknownJob_ReturnsError(t *testing.T) {
	w := NewForTest(testConfig())
	err := w.RecordExitCode("nonexistent", 0)
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetExitCode_UnknownJob_ReturnsFalse(t *testing.T) {
	w := NewForTest(testConfig())
	_, ok := w.GetExitCode("nonexistent")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestRecordExitCode_OverwritesPrevious(t *testing.T) {
	w := NewForTest(testConfig())
	_ = w.RecordExitCode("backup", 0)
	_ = w.RecordExitCode("backup", 2)
	e, ok := w.GetExitCode("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Code != 2 {
		t.Errorf("expected code 2, got %d", e.Code)
	}
}
