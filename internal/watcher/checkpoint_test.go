package watcher

import (
	"testing"
)

func TestRecordCheckpoint_StoresEntry(t *testing.T) {
	w := NewForTest(testConfig())
	if err := w.RecordCheckpoint("backup", "files-copied", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, ok := w.GetCheckpoints("backup")
	if !ok {
		t.Fatal("expected checkpoints to exist")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Name != "files-copied" {
		t.Errorf("expected name %q, got %q", "files-copied", entries[0].Name)
	}
}

func TestRecordCheckpoint_UnknownJob_ReturnsError(t *testing.T) {
	w := NewForTest(testConfig())
	if err := w.RecordCheckpoint("nonexistent", "step", nil); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetCheckpoints_UnknownJob_ReturnsFalse(t *testing.T) {
	w := NewForTest(testConfig())
	_, ok := w.GetCheckpoints("nonexistent")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestClearCheckpoints_RemovesEntries(t *testing.T) {
	w := NewForTest(testConfig())
	_ = w.RecordCheckpoint("backup", "step-1", nil)
	_ = w.RecordCheckpoint("backup", "step-2", nil)
	if err := w.ClearCheckpoints("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, ok := w.GetCheckpoints("backup")
	if !ok {
		t.Fatal("expected job to still exist")
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(entries))
	}
}

func TestRecordCheckpoint_StoresMeta(t *testing.T) {
	w := NewForTest(testConfig())
	meta := map[string]string{"files": "42"}
	_ = w.RecordCheckpoint("backup", "counted", meta)
	entries, _ := w.GetCheckpoints("backup")
	if entries[0].Meta["files"] != "42" {
		t.Errorf("expected meta files=42, got %v", entries[0].Meta)
	}
}

func TestGetCheckpoints_ReturnsCopy(t *testing.T) {
	w := NewForTest(testConfig())
	_ = w.RecordCheckpoint("backup", "step", nil)
	a, _ := w.GetCheckpoints("backup")
	a[0].Name = "mutated"
	b, _ := w.GetCheckpoints("backup")
	if b[0].Name == "mutated" {
		t.Error("GetCheckpoints should return a copy")
	}
}
