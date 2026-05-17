package watcher

import (
	"fmt"
	"sync"
	"time"
)

// RetryEntry records a single retry attempt for a job.
type RetryEntry struct {
	Attempt   int
	At        time.Time
	Reason    string
}

// RetrySummary holds the current retry state for a job.
type RetrySummary struct {
	JobName     string
	Attempts    int
	LastAttempt time.Time
	LastReason  string
}

type retryStore struct {
	mu      sync.RWMutex
	entries map[string][]RetryEntry
}

func newRetryStore() *retryStore {
	return &retryStore{
		entries: make(map[string][]RetryEntry),
	}
}

func (s *retryStore) recordRetry(job, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.entries[job]; !ok {
		return fmt.Errorf("unknown job: %s", job)
	}

	attempt := len(s.entries[job]) + 1
	s.entries[job] = append(s.entries[job], RetryEntry{
		Attempt: attempt,
		At:      time.Now(),
		Reason:  reason,
	})
	return nil
}

func (s *retryStore) getRetrySummary(job string) (RetrySummary, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, ok := s.entries[job]
	if !ok {
		return RetrySummary{}, false
	}

	summary := RetrySummary{
		JobName:  job,
		Attempts: len(entries),
	}
	if len(entries) > 0 {
		last := entries[len(entries)-1]
		summary.LastAttempt = last.At
		summary.LastReason = last.Reason
	}
	return summary, true
}

func (s *retryStore) clearRetries(job string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.entries[job]; !ok {
		return fmt.Errorf("unknown job: %s", job)
	}
	s.entries[job] = nil
	return nil
}
