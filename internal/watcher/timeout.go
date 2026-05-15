package watcher

import (
	"fmt"
	"sync"
	"time"
)

// TimeoutEntry records when a job was started and its configured max duration.
type TimeoutEntry struct {
	StartedAt   time.Time
	MaxDuration time.Duration
}

type timeoutStore struct {
	mu      sync.RWMutex
	entries map[string]TimeoutEntry
}

func newTimeoutStore() *timeoutStore {
	return &timeoutStore{
		entries: make(map[string]TimeoutEntry),
	}
}

// RecordTimeout registers a running job with its maximum allowed duration.
func (w *Watcher) RecordTimeout(name string, max time.Duration) error {
	w.mu.RLock()
	_, known := w.jobs[name]
	w.mu.RUnlock()
	if !known {
		return fmt.Errorf("unknown job: %s", name)
	}

	w.timeouts.mu.Lock()
	defer w.timeouts.mu.Unlock()
	w.timeouts.entries[name] = TimeoutEntry{
		StartedAt:   w.clock.Now(),
		MaxDuration: max,
	}
	return nil
}

// ClearTimeout removes the timeout entry for a job (called on finish).
func (w *Watcher) ClearTimeout(name string) {
	w.timeouts.mu.Lock()
	defer w.timeouts.mu.Unlock()
	delete(w.timeouts.entries, name)
}

// CheckTimeouts returns the names of jobs that have exceeded their max duration.
func (w *Watcher) CheckTimeouts() []string {
	now := w.clock.Now()
	w.timeouts.mu.RLock()
	defer w.timeouts.mu.RUnlock()

	var overdue []string
	for name, entry := range w.timeouts.entries {
		if now.Sub(entry.StartedAt) > entry.MaxDuration {
			overdue = append(overdue, name)
		}
	}
	return overdue
}
