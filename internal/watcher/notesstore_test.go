package watcher

import (
	"testing"
)

func TestAddNote_StoresEntry(t *testing.T) {
	s := newNotesStore([]string{"backup"})
	entry, err := s.addNote("backup", "manual check required")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Text != "manual check required" {
		t.Errorf("expected text %q, got %q", "manual check required", entry.Text)
	}
	if entry.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestAddNote_UnknownJob_ReturnsError(t *testing.T) {
	s := newNotesStore([]string{})
	_, err := s.addNote("ghost", "hello")
	if err == nil {
		t.Error("expected error for unknown job")
	}
}

func TestGetNotes_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newNotesStore([]string{})
	_, ok := s.getNotes("ghost")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestGetNotes_ReturnsCopy(t *testing.T) {
	s := newNotesStore([]string{"backup"})
	_, _ = s.addNote("backup", "note one")
	_, _ = s.addNote("backup", "note two")
	entries, ok := s.getNotes("backup")
	if !ok {
		t.Fatal("expected ok")
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestDeleteNote_RemovesEntry(t *testing.T) {
	s := newNotesStore([]string{"backup"})
	entry, _ := s.addNote("backup", "to be deleted")
	err := s.deleteNote("backup", entry.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := s.getNotes("backup")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after delete, got %d", len(entries))
	}
}

func TestDeleteNote_UnknownID_ReturnsError(t *testing.T) {
	s := newNotesStore([]string{"backup"})
	err := s.deleteNote("backup", "nonexistent-id")
	if err == nil {
		t.Error("expected error for unknown note ID")
	}
}

func TestDeleteNote_UnknownJob_ReturnsError(t *testing.T) {
	s := newNotesStore([]string{})
	err := s.deleteNote("ghost", "any-id")
	if err == nil {
		t.Error("expected error for unknown job")
	}
}
