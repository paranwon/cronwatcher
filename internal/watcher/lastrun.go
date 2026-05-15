package watcher

import (
	"fmt"
	"sync"
	"time"
)

// LastRunEntry holds information about the most recent completed run of a job.
type LastRunEntry struct {
	FinishedAt time.Time
	Duration   time.Duration
	Success    bool
	Error      string
}

type lastRunStore struct {
	mu      sync.RWMutex
	entries map[string]*LastRunEntry
}

func newLastRunStore() *lastRunStore {
	return &lastRunStore{
		entries: make(map[string]*LastRunEntry),
	}
}

func (s *lastRunStore) record(name string, entry LastRunEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[name] = &entry
}

func (s *lastRunStore) get(name string) (LastRunEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[name]
	if !ok {
		return LastRunEntry{}, false
	}
	return *e, true
}

// RecordLastRun stores the result of the most recent completed run for a job.
func (w *Watcher) RecordLastRun(name string, entry LastRunEntry) error {
	if _, ok := w.jobs[name]; !ok {
		return fmt.Errorf("unknown job: %s", name)
	}
	w.lastRun.record(name, entry)
	return nil
}

// GetLastRun returns the most recent run entry for the named job.
func (w *Watcher) GetLastRun(name string) (LastRunEntry, bool) {
	return w.lastRun.get(name)
}
