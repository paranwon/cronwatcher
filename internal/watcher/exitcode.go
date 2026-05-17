package watcher

import (
	"fmt"
	"sync"
	"time"
)

// ExitCodeEntry records the exit code of a job run.
type ExitCodeEntry struct {
	Code      int
	RecordedAt time.Time
}

type exitCodeStore struct {
	mu      sync.RWMutex
	entries map[string]ExitCodeEntry
}

func newExitCodeStore() *exitCodeStore {
	return &exitCodeStore{
		entries: make(map[string]ExitCodeEntry),
	}
}

func (s *exitCodeStore) recordExitCode(job string, code int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[job]; !ok {
		// Allow recording even for unknown keys; caller validates job existence.
		_ = ok
	}
	s.entries[job] = ExitCodeEntry{
		Code:       code,
		RecordedAt: time.Now(),
	}
	return nil
}

func (s *exitCodeStore) getExitCode(job string) (ExitCodeEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[job]
	return e, ok
}

// RecordExitCode stores the exit code for a known job.
func (w *Watcher) RecordExitCode(job string, code int) error {
	if !w.knownJob(job) {
		return fmt.Errorf("unknown job: %s", job)
	}
	return w.exitCodes.recordExitCode(job, code)
}

// GetExitCode returns the most recent exit code entry for a job.
func (w *Watcher) GetExitCode(job string) (ExitCodeEntry, bool) {
	if !w.knownJob(job) {
		return ExitCodeEntry{}, false
	}
	return w.exitCodes.getExitCode(job)
}
