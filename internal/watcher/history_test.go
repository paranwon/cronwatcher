package watcher

import (
	"errors"
	"testing"
	"time"
)

func TestGetHistory_ReturnsEntriesForKnownJob(t *testing.T) {
	w := NewForTest(testConfig())
	w.RecordStart("backup")
	time.Sleep(2 * time.Millisecond)
	w.RecordFinish("backup", nil)

	entries, ok := w.GetHistory("backup", 10)
	if !ok {
		t.Fatal("expected job to be found")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Duration == "" {
		t.Error("expected non-empty duration")
	}
}

func TestGetHistory_UnknownJob_ReturnsFalse(t *testing.T) {
	w := NewForTest(testConfig())
	_, ok := w.GetHistory("nonexistent", 10)
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetHistory_RespectsLimit(t *testing.T) {
	w := NewForTest(testConfig())
	for i := 0; i < 5; i++ {
		w.RecordStart("backup")
		w.RecordFinish("backup", nil)
	}

	entries, ok := w.GetHistory("backup", 3)
	if !ok {
		t.Fatal("expected job to be found")
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestGetHistory_RecordsError(t *testing.T) {
	w := NewForTest(testConfig())
	w.RecordStart("backup")
	w.RecordFinish("backup", errors.New("exit status 1"))

	entries, _ := w.GetHistory("backup", 10)
	if entries[0].Error != "exit status 1" {
		t.Errorf("expected error string, got %q", entries[0].Error)
	}
}
