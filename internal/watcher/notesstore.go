package watcher

import (
	"fmt"
	"sync"
	"time"
)

// NoteEntry represents a single operator note attached to a job.
type NoteEntry struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type notesStore struct {
	mu    sync.RWMutex
	notes map[string][]NoteEntry
}

func newNotesStore(jobs []string) *notesStore {
	m := make(map[string][]NoteEntry, len(jobs))
	for _, j := range jobs {
		m[j] = []NoteEntry{}
	}
	return &notesStore{notes: m}
}

func (s *notesStore) addNote(job, text string) (NoteEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notes[job]; !ok {
		return NoteEntry{}, fmt.Errorf("unknown job: %s", job)
	}
	entry := NoteEntry{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Text:      text,
		CreatedAt: time.Now().UTC(),
	}
	s.notes[job] = append(s.notes[job], entry)
	return entry, nil
}

func (s *notesStore) getNotes(job string) ([]NoteEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, ok := s.notes[job]
	if !ok {
		return nil, false
	}
	copy := make([]NoteEntry, len(entries))
	for i, e := range entries {
		copy[i] = e
	}
	return copy, true
}

func (s *notesStore) deleteNote(job, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, ok := s.notes[job]
	if !ok {
		return fmt.Errorf("unknown job: %s", job)
	}
	for i, e := range entries {
		if e.ID == id {
			s.notes[job] = append(entries[:i], entries[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("note not found: %s", id)
}
