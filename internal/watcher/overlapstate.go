package watcher

import (
	"fmt"
	"sync"
	"time"
)

// OverlapEntry records when a job was detected running concurrently with a previous instance.
type OverlapEntry struct {
	JobName    string
	DetectedAt time.Time
	RunSince   time.Time
}

type overlapStore struct {
	mu      sync.RWMutex
	entries map[string][]OverlapEntry
}

func newOverlapStore() *overlapStore {
	return &overlapStore{
		entries: make(map[string][]OverlapEntry),
	}
}

// RecordOverlap records a concurrent-execution overlap event for a known job.
func (w *Watcher) RecordOverlap(jobName string, runningSince time.Time) error {
	w.overlap.mu.Lock()
	defer w.overlap.mu.Unlock()

	if _, ok := w.jobs[jobName]; !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}

	entry := OverlapEntry{
		JobName:    jobName,
		DetectedAt: time.Now(),
		RunSince:   runningSince,
	}
	w.overlap.entries[jobName] = append(w.overlap.entries[jobName], entry)
	return nil
}

// GetOverlaps returns all recorded overlap events for a job.
// Returns false if the job is unknown.
func (w *Watcher) GetOverlaps(jobName string) ([]OverlapEntry, bool) {
	w.overlap.mu.RLock()
	defer w.overlap.mu.RUnlock()

	if _, ok := w.jobs[jobName]; !ok {
		return nil, false
	}

	src := w.overlap.entries[jobName]
	out := make([]OverlapEntry, len(src))
	copy(out, src)
	return out, true
}

// ClearOverlaps removes all overlap entries for a job.
func (w *Watcher) ClearOverlaps(jobName string) error {
	w.overlap.mu.Lock()
	defer w.overlap.mu.Unlock()

	if _, ok := w.jobs[jobName]; !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}

	delete(w.overlap.entries, jobName)
	return nil
}
