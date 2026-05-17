package watcher

import (
	"fmt"
	"sync"
	"time"
)

// BlockedEntry records when a job was blocked by a dependency or lock.
type BlockedEntry struct {
	JobName   string
	BlockedBy string
	At        time.Time
	Reason    string
}

type blockedStore struct {
	mu      sync.RWMutex
	entries map[string][]BlockedEntry
}

func newBlockedStore(jobs []string) *blockedStore {
	m := make(map[string][]BlockedEntry, len(jobs))
	for _, j := range jobs {
		m[j] = []BlockedEntry{}
	}
	return &blockedStore{entries: m}
}

func (s *blockedStore) recordBlocked(jobName, blockedBy, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[jobName]; !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}
	s.entries[jobName] = append(s.entries[jobName], BlockedEntry{
		JobName:   jobName,
		BlockedBy: blockedBy,
		At:        time.Now(),
		Reason:    reason,
	})
	return nil
}

func (s *blockedStore) getBlocked(jobName string) ([]BlockedEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, ok := s.entries[jobName]
	if !ok {
		return nil, false
	}
	copy := make([]BlockedEntry, len(entries))
	for i, e := range entries {
		copy[i] = e
	}
	return copy, true
}

func (s *blockedStore) clearBlocked(jobName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[jobName]; !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}
	s.entries[jobName] = []BlockedEntry{}
	return nil
}
