package watcher

import "fmt"

// PauseJob marks a job as paused so missed-job alerts are suppressed.
func (w *Watcher) PauseJob(name string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	state, ok := w.jobs[name]
	if !ok {
		return fmt.Errorf("job %q not found", name)
	}
	state.Paused = true
	w.jobs[name] = state
	return nil
}

// ResumeJob clears the paused flag for a job.
func (w *Watcher) ResumeJob(name string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	state, ok := w.jobs[name]
	if !ok {
		return fmt.Errorf("job %q not found", name)
	}
	state.Paused = false
	w.jobs[name] = state
	return nil
}

// IsPaused returns true if the named job is currently paused.
func (w *Watcher) IsPaused(name string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	state, ok := w.jobs[name]
	return ok && state.Paused
}
