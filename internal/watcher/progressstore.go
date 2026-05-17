package watcher

import (
	"fmt"
	"sync"
	"time"
)

// ProgressEntry holds a progress snapshot for a running job.
type ProgressEntry struct {
	Step      int       `json:"step"`
	Total     int       `json:"total"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Percent returns the completion percentage, or 0 if total is zero.
func (p ProgressEntry) Percent() float64 {
	if p.Total == 0 {
		return 0
	}
	return float64(p.Step) / float64(p.Total) * 100
}

type progressStore struct {
	mu      sync.RWMutex
	entries map[string]ProgressEntry
}

func newProgressStore(jobs []string) *progressStore {
	entries := make(map[string]ProgressEntry, len(jobs))
	for _, j := range jobs {
		entries[j] = ProgressEntry{}
	}
	return &progressStore{entries: entries}
}

func (s *progressStore) recordProgress(job string, step, total int, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[job]; !ok {
		return fmt.Errorf("unknown job: %s", job)
	}
	s.entries[job] = ProgressEntry{
		Step:      step,
		Total:     total,
		Message:   message,
		UpdatedAt: time.Now(),
	}
	return nil
}

func (s *progressStore) getProgress(job string) (ProgressEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[job]
	return e, ok
}

func (s *progressStore) clearProgress(job string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[job]; !ok {
		return fmt.Errorf("unknown job: %s", job)
	}
	s.entries[job] = ProgressEntry{}
	return nil
}
