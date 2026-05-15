package watcher

import (
	"fmt"
	"sync"
	"time"
)

// NextRunEntry holds the predicted next run time for a job.
type NextRunEntry struct {
	JobName     string
	NextRun     time.Time
	Schedule    string
	LastUpdated time.Time
}

type nextRunStore struct {
	mu      sync.RWMutex
	entries map[string]NextRunEntry
}

func newNextRunStore() *nextRunStore {
	return &nextRunStore{
		entries: make(map[string]NextRunEntry),
	}
}

// SetNextRun stores the predicted next run time for a job.
func (w *Watcher) SetNextRun(jobName string, next time.Time, schedule string) error {
	w.mu.RLock()
	_, known := w.jobs[jobName]
	w.mu.RUnlock()
	if !known {
		return fmt.Errorf("unknown job: %s", jobName)
	}

	w.nextRun.mu.Lock()
	defer w.nextRun.mu.Unlock()
	w.nextRun.entries[jobName] = NextRunEntry{
		JobName:     jobName,
		NextRun:     next,
		Schedule:    schedule,
		LastUpdated: time.Now(),
	}
	return nil
}

// GetNextRun returns the predicted next run entry for a job.
func (w *Watcher) GetNextRun(jobName string) (NextRunEntry, bool) {
	w.nextRun.mu.RLock()
	defer w.nextRun.mu.RUnlock()
	e, ok := w.nextRun.entries[jobName]
	return e, ok
}

// OverdueJobs returns jobs whose predicted next run is in the past.
func (w *Watcher) OverdueJobs(now time.Time) []NextRunEntry {
	w.nextRun.mu.RLock()
	defer w.nextRun.mu.RUnlock()

	var overdue []NextRunEntry
	for _, e := range w.nextRun.entries {
		if !e.NextRun.IsZero() && now.After(e.NextRun) {
			overdue = append(overdue, e)
		}
	}
	return overdue
}
