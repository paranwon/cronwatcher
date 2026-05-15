package watcher

import "fmt"

// SetAnnotations stores arbitrary key-value annotations for a job.
// Annotations are merged with any existing ones.
func (w *Watcher) SetAnnotations(jobName string, annotations map[string]string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	state, ok := w.jobs[jobName]
	if !ok {
		return fmt.Errorf("job %q not found", jobName)
	}

	if state.Annotations == nil {
		state.Annotations = make(map[string]string, len(annotations))
	}

	for k, v := range annotations {
		state.Annotations[k] = v
	}

	w.jobs[jobName] = state
	return nil
}

// GetAnnotations returns a copy of the annotations for the given job.
// Returns false if the job is not known.
func (w *Watcher) GetAnnotations(jobName string) (map[string]string, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	state, ok := w.jobs[jobName]
	if !ok {
		return nil, false
	}

	copy := make(map[string]string, len(state.Annotations))
	for k, v := range state.Annotations {
		copy[k] = v
	}
	return copy, true
}

// DeleteAnnotation removes a single annotation key from a job.
func (w *Watcher) DeleteAnnotation(jobName, key string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	state, ok := w.jobs[jobName]
	if !ok {
		return fmt.Errorf("job %q not found", jobName)
	}

	delete(state.Annotations, key)
	w.jobs[jobName] = state
	return nil
}
