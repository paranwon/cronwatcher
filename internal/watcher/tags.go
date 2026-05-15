package watcher

import "fmt"

// SetTags assigns metadata tags to a registered job.
// Tags are arbitrary key-value pairs (e.g. team, env, owner).
func (w *Watcher) SetTags(jobName string, tags map[string]string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	state, ok := w.jobs[jobName]
	if !ok {
		return fmt.Errorf("watcher: unknown job %q", jobName)
	}

	if state.Tags == nil {
		state.Tags = make(map[string]string, len(tags))
	}
	for k, v := range tags {
		state.Tags[k] = v
	}
	w.jobs[jobName] = state
	return nil
}

// GetTags returns the metadata tags for a registered job.
// Returns false if the job is not known.
func (w *Watcher) GetTags(jobName string) (map[string]string, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	state, ok := w.jobs[jobName]
	if !ok {
		return nil, false
	}

	// Return a copy to avoid external mutation.
	copy := make(map[string]string, len(state.Tags))
	for k, v := range state.Tags {
		copy[k] = v
	}
	return copy, true
}
