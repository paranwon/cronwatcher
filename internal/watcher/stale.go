package watcher

import (
	"fmt"
	"time"
)

// StaleEntry describes a job that has not been seen within its expected window.
type StaleEntry struct {
	JobName  string
	LastSeen time.Time
	StaleSince time.Duration
}

// MarkStale returns all jobs whose last recorded finish time exceeds the
// configured stale threshold. A job is considered stale when it has not
// completed successfully within 2× its expected interval.
func (w *Watcher) MarkStale(now time.Time) ([]StaleEntry, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var stale []StaleEntry

	for _, job := range w.config.Jobs {
		state, ok := w.jobs[job.Name]
		if !ok {
			continue
		}

		// Skip jobs that are currently running or paused.
		if state.Running || state.Paused {
			continue
		}

		if state.LastSeen.IsZero() {
			continue
		}

		threshold := job.Interval * 2
		age := now.Sub(state.LastSeen)
		if age > threshold {
			stale = append(stale, StaleEntry{
				JobName:    job.Name,
				LastSeen:   state.LastSeen,
				StaleSince: age - threshold,
			})
		}
	}

	return stale, nil
}

// StaleError is returned when a job name is not registered.
type StaleError struct {
	JobName string
}

func (e *StaleError) Error() string {
	return fmt.Sprintf("stale check: unknown job %q", e.JobName)
}
