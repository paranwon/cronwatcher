package watcher

import "fmt"

// SetLabels sets or merges key-value labels on a job.
func (w *Watcher) SetLabels(jobName string, labels map[string]string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	state, ok := w.jobs[jobName]
	if !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}

	if state.Labels == nil {
		state.Labels = make(map[string]string)
	}
	for k, v := range labels {
		state.Labels[k] = v
	}
	w.jobs[jobName] = state
	return nil
}

// GetLabels returns a copy of the labels for a job.
func (w *Watcher) GetLabels(jobName string) (map[string]string, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	state, ok := w.jobs[jobName]
	if !ok {
		return nil, false
	}

	copy := make(map[string]string, len(state.Labels))
	for k, v := range state.Labels {
		copy[k] = v
	}
	return copy, true
}

// DeleteLabel removes a single label key from a job.
func (w *Watcher) DeleteLabel(jobName, key string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	state, ok := w.jobs[jobName]
	if !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}
	delete(state.Labels, key)
	w.jobs[jobName] = state
	return nil
}
