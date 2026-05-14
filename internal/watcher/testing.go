package watcher

import (
	"time"
)

// NewForTest returns a Watcher pre-configured for unit tests.
func NewForTest() *Watcher {
	return &Watcher{
		jobs:     make(map[string]*JobState),
		statuses: make(map[string]Status),
	}
}

// InjectMissed inserts a missed job status directly into the watcher for testing.
func (w *Watcher) InjectMissed(name string, lastRun time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.statuses[name] = Status{
		Name:    name,
		Missed:  true,
		LastRun: lastRun,
	}
}

// InjectLongRunning inserts a long-running job status directly into the watcher for testing.
func (w *Watcher) InjectLongRunning(name string, duration time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.statuses[name] = Status{
		Name:             name,
		Running:          true,
		Duration:         duration,
		ExceedsThreshold: true,
	}
}
