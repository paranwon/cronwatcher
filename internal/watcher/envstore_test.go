package watcher

import (
	"testing"
)

func TestRecordEnv_StoresEntry(t *testing.T) {
	w := NewForTest([]string{"backup"})

	err := w.RecordEnv("backup", map[string]string{"HOME": "/root", "PATH": "/usr/bin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry, ok := w.GetEnv("backup")
	if !ok {
		t.Fatal("expected env entry to exist")
	}
	if entry.Vars["HOME"] != "/root" {
		t.Errorf("expected HOME=/root, got %s", entry.Vars["HOME"])
	}
	if entry.Vars["PATH"] != "/usr/bin" {
		t.Errorf("expected PATH=/usr/bin, got %s", entry.Vars["PATH"])
	}
}

func TestRecordEnv_UnknownJob_ReturnsError(t *testing.T) {
	w := NewForTest([]string{"backup"})

	err := w.RecordEnv("nonexistent", map[string]string{"KEY": "val"})
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetEnv_UnknownJob_ReturnsFalse(t *testing.T) {
	w := NewForTest([]string{"backup"})

	_, ok := w.GetEnv("nonexistent")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetEnv_ReturnsCopy(t *testing.T) {
	w := NewForTest([]string{"backup"})

	_ = w.RecordEnv("backup", map[string]string{"KEY": "original"})

	entry, _ := w.GetEnv("backup")
	entry.Vars["KEY"] = "mutated"

	entry2, _ := w.GetEnv("backup")
	if entry2.Vars["KEY"] != "original" {
		t.Errorf("expected original value, got %s", entry2.Vars["KEY"])
	}
}

func TestClearEnv_RemovesEntry(t *testing.T) {
	w := NewForTest([]string{"backup"})

	_ = w.RecordEnv("backup", map[string]string{"KEY": "val"})
	err := w.ClearEnv("backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := w.GetEnv("backup")
	if ok {
		t.Fatal("expected env entry to be cleared")
	}
}

func TestClearEnv_UnknownJob_ReturnsError(t *testing.T) {
	w := NewForTest([]string{"backup"})

	err := w.ClearEnv("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}
