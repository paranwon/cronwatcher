package watcher

import (
	"fmt"
	"sync"
	"time"
)

// OutputEntry holds a captured output snapshot for a job run.
type OutputEntry struct {
	Timestamp time.Time
	Stdout    string
	Stderr    string
	Truncated bool
}

const maxOutputBytes = 4096

type outputStore struct {
	mu      sync.RWMutex
	entries map[string]*OutputEntry
}

func newOutputStore() *outputStore {
	return &outputStore{
		entries: make(map[string]*OutputEntry),
	}
}

func (s *outputStore) recordOutput(job, stdout, stderr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.entries[job]; !ok {
		return fmt.Errorf("unknown job: %s", job)
	}

	truncated := false
	if len(stdout) > maxOutputBytes {
		stdout = stdout[:maxOutputBytes]
		truncated = true
	}
	if len(stderr) > maxOutputBytes {
		stderr = stderr[:maxOutputBytes]
		truncated = true
	}

	s.entries[job] = &OutputEntry{
		Timestamp: time.Now(),
		Stdout:    stdout,
		Stderr:    stderr,
		Truncated: truncated,
	}
	return nil
}

func (s *outputStore) getOutput(job string) (OutputEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.entries[job]
	if !ok || e == nil {
		return OutputEntry{}, false
	}
	return *e, true
}

func (s *outputStore) initJob(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[job]; !ok {
		s.entries[job] = nil
	}
}
