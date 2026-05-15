package watcher

import (
	"fmt"
	"sync"
	"time"
)

// SLAEntry records whether a job met its SLA window.
type SLAEntry struct {
	JobName   string
	Deadline  time.Time
	MetSLA    bool
	CheckedAt time.Time
}

type slaStore struct {
	mu      sync.RWMutex
	entries map[string]SLAEntry
}

func newSLAStore() *slaStore {
	return &slaStore{
		entries: make(map[string]SLAEntry),
	}
}

// RecordSLA records the SLA evaluation result for a job.
func (s *slaStore) RecordSLA(name string, deadline time.Time, met bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[name]; !ok {
		// allow recording even for unknown jobs; caller validates existence
	}
	s.entries[name] = SLAEntry{
		JobName:   name,
		Deadline:  deadline,
		MetSLA:    met,
		CheckedAt: time.Now(),
	}
	return nil
}

// GetSLA returns the latest SLA entry for a job.
func (s *slaStore) GetSLA(name string) (SLAEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[name]
	return e, ok
}

// CheckSLAViolations returns jobs that missed their SLA deadline.
func (s *slaStore) CheckSLAViolations() []SLAEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var violations []SLAEntry
	for _, e := range s.entries {
		if !e.MetSLA {
			violations = append(violations, e)
		}
	}
	return violations
}

// slaKey is used for map access validation.
func slaKey(name string) string { return fmt.Sprintf("sla:%s", name) }
