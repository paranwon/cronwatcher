package watcher

import (
	"sync"
	"time"
)

// AvgDurationEntry holds rolling average duration stats for a job.
type AvgDurationEntry struct {
	Count   int64
	Total   time.Duration
	Average time.Duration
}

type avgDurationStore struct {
	mu      sync.RWMutex
	entries map[string]*AvgDurationEntry
}

func newAvgDurationStore(jobs []string) *avgDurationStore {
	s := &avgDurationStore{
		entries: make(map[string]*AvgDurationEntry, len(jobs)),
	}
	for _, name := range jobs {
		s.entries[name] = &AvgDurationEntry{}
	}
	return s
}

// RecordDuration adds a completed job duration to the rolling stats.
func (w *Watcher) RecordDuration(name string, d time.Duration) error {
	w.avgDuration.mu.Lock()
	defer w.avgDuration.mu.Unlock()

	e, ok := w.avgDuration.entries[name]
	if !ok {
		return errUnknownJob(name)
	}
	e.Count++
	e.Total += d
	e.Average = time.Duration(int64(e.Total) / e.Count)
	return nil
}

// GetAvgDuration returns the average duration entry for a job.
func (w *Watcher) GetAvgDuration(name string) (AvgDurationEntry, bool) {
	w.avgDuration.mu.RLock()
	defer w.avgDuration.mu.RUnlock()

	e, ok := w.avgDuration.entries[name]
	if !ok {
		return AvgDurationEntry{}, false
	}
	return *e, true
}
