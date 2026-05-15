package watcher

import "time"

// HistoryEntry records a single completed execution of a job.
type HistoryEntry struct {
	StartedAt  time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Duration   string        `json:"duration"`
	Error      string        `json:"error,omitempty"`
}

// RecordHistory appends a completed run to the job's history ring buffer.
func (w *Watcher) RecordHistory(name string, start, finish time.Time, runErr error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	entry := HistoryEntry{
		StartedAt:  start,
		FinishedAt: finish,
		Duration:   formatDuration(finish.Sub(start)),
	}
	if runErr != nil {
		entry.Error = runErr.Error()
	}

	w.history[name] = append(w.history[name], entry)
}

// GetHistory returns up to limit recent history entries for the named job.
// Returns false if the job is unknown.
func (w *Watcher) GetHistory(name string, limit int) ([]HistoryEntry, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if _, known := w.jobs[name]; !known {
		return nil, false
	}

	entries := w.history[name]
	if len(entries) <= limit {
		return entries, true
	}
	return entries[len(entries)-limit:], true
}
