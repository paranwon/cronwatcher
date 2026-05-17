package watcher

import (
	"fmt"
	"sync"
	"time"
)

// CustomStatus holds an arbitrary status message posted by a job.
type CustomStatus struct {
	Message   string
	Severity  string // "info", "warn", "error"
	RecordedAt time.Time
}

type customStatusStore struct {
	mu      sync.RWMutex
	entries map[string]CustomStatus
}

func newCustomStatusStore() *customStatusStore {
	return &customStatusStore{
		entries: make(map[string]CustomStatus),
	}
}

// RecordCustomStatus stores a custom status message for a known job.
func (w *Watcher) RecordCustomStatus(job, message, severity string) error {
	w.mu.RLock()
	_, known := w.jobs[job]
	w.mu.RUnlock()
	if !known {
		return fmt.Errorf("unknown job: %s", job)
	}
	w.customStatus.mu.Lock()
	defer w.customStatus.mu.Unlock()
	w.customStatus.entries[job] = CustomStatus{
		Message:    message,
		Severity:   severity,
		RecordedAt: time.Now(),
	}
	return nil
}

// GetCustomStatus retrieves the custom status for a job.
func (w *Watcher) GetCustomStatus(job string) (CustomStatus, bool) {
	w.customStatus.mu.RLock()
	defer w.customStatus.mu.RUnlock()
	s, ok := w.customStatus.entries[job]
	return s, ok
}

// ClearCustomStatus removes the custom status entry for a job.
func (w *Watcher) ClearCustomStatus(job string) error {
	w.mu.RLock()
	_, known := w.jobs[job]
	w.mu.RUnlock()
	if !known {
		return fmt.Errorf("unknown job: %s", job)
	}
	w.customStatus.mu.Lock()
	defer w.customStatus.mu.Unlock()
	delete(w.customStatus.entries, job)
	return nil
}
