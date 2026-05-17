package watcher

import (
	"fmt"
	"sync"
	"time"
)

// TriggerEntry records a manual trigger event for a job.
type TriggerEntry struct {
	TriggeredAt time.Time `json:"triggered_at"`
	TriggeredBy string    `json:"triggered_by"`
	Reason      string    `json:"reason,omitempty"`
}

type triggerStore struct {
	mu      sync.RWMutex
	entries map[string]*TriggerEntry
}

func newTriggerStore() *triggerStore {
	return &triggerStore{
		entries: make(map[string]*TriggerEntry),
	}
}

func (s *triggerStore) recordTrigger(job, triggeredBy, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[job]; !ok {
		return fmt.Errorf("unknown job: %s", job)
	}
	s.entries[job] = &TriggerEntry{
		TriggeredAt: time.Now(),
		TriggeredBy: triggeredBy,
		Reason:      reason,
	}
	return nil
}

func (s *triggerStore) getTrigger(job string) (TriggerEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[job]
	if !ok || e == nil {
		return TriggerEntry{}, false
	}
	return *e, true
}

func (s *triggerStore) clearTrigger(job string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[job]; !ok {
		return fmt.Errorf("unknown job: %s", job)
	}
	s.entries[job] = nil
	return nil
}

func (s *triggerStore) register(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[job] = nil
}
