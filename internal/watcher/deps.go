package watcher

import "time"

// Dependency tracks inter-job dependencies, recording which jobs
// must succeed before a dependent job is considered healthy.
type Dependency struct {
	RequiredJob string
	MaxStaleness time.Duration
}

// SetDependencies stores the dependency list for the given job.
func (w *Watcher) SetDependencies(jobName string, deps []Dependency) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.jobs[jobName]; !ok {
		return ErrUnknownJob
	}
	w.deps[jobName] = deps
	return nil
}

// GetDependencies returns the dependency list for the given job.
// Returns false if the job is unknown.
func (w *Watcher) GetDependencies(jobName string) ([]Dependency, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if _, ok := w.jobs[jobName]; !ok {
		return nil, false
	}
	deps, ok := w.deps[jobName]
	if !ok {
		return []Dependency{}, true
	}
	out := make([]Dependency, len(deps))
	copy(out, deps)
	return out, true
}

// CheckDependencies returns the names of required jobs whose last
// successful finish is older than the configured MaxStaleness.
// An empty slice means all dependencies are satisfied.
func (w *Watcher) CheckDependencies(jobName string) []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	deps, ok := w.deps[jobName]
	if !ok {
		return nil
	}

	now := time.Now()
	var unsatisfied []string
	for _, dep := range deps {
		st, exists := w.jobs[dep.RequiredJob]
		if !exists {
			unsatisfied = append(unsatisfied, dep.RequiredJob)
			continue
		}
		if st.LastSuccess.IsZero() || now.Sub(st.LastSuccess) > dep.MaxStaleness {
			unsatisfied = append(unsatisfied, dep.RequiredJob)
		}
	}
	return unsatisfied
}
