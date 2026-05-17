package watcher

import (
	"testing"
)

func TestRecordProgress_StoresEntry(t *testing.T) {
	s := newProgressStore([]string{"backup"})
	if err := s.recordProgress("backup", 3, 10, "processing"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.getProgress("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Step != 3 || e.Total != 10 || e.Message != "processing" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestRecordProgress_UnknownJob_ReturnsError(t *testing.T) {
	s := newProgressStore([]string{})
	if err := s.recordProgress("ghost", 1, 5, "hi"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetProgress_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newProgressStore([]string{})
	_, ok := s.getProgress("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestProgressEntry_Percent_CalculatesCorrectly(t *testing.T) {
	e := ProgressEntry{Step: 1, Total: 4}
	if e.Percent() != 25.0 {
		t.Errorf("expected 25.0, got %f", e.Percent())
	}
}

func TestProgressEntry_Percent_ZeroTotal_ReturnsZero(t *testing.T) {
	e := ProgressEntry{Step: 5, Total: 0}
	if e.Percent() != 0 {
		t.Errorf("expected 0, got %f", e.Percent())
	}
}

func TestClearProgress_ResetsEntry(t *testing.T) {
	s := newProgressStore([]string{"sync"})
	_ = s.recordProgress("sync", 7, 10, "almost done")
	if err := s.clearProgress("sync"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.getProgress("sync")
	if !ok {
		t.Fatal("expected entry to still exist after clear")
	}
	if e.Step != 0 || e.Total != 0 || e.Message != "" {
		t.Errorf("expected zeroed entry after clear, got %+v", e)
	}
}

func TestClearProgress_UnknownJob_ReturnsError(t *testing.T) {
	s := newProgressStore([]string{})
	if err := s.clearProgress("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}
