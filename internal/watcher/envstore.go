package watcher

import (
	"fmt"
	"sync"
)

// EnvEntry holds environment variable snapshot for a job run.
type EnvEntry struct {
	Vars map[string]string
}

type envStore struct {
	mu      sync.RWMutex
	entries map[string]EnvEntry
}

func newEnvStore() *envStore {
	return &envStore{
		entries: make(map[string]EnvEntry),
	}
}

// RecordEnv stores the environment variable snapshot for a known job.
func (w *Watcher) RecordEnv(jobName string, vars map[string]string) error {
	w.env.mu.Lock()
	defer w.env.mu.Unlock()

	if _, ok := w.jobs[jobName]; !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}

	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}

	w.env.entries[jobName] = EnvEntry{Vars: copy}
	return nil
}

// GetEnv returns the last recorded environment snapshot for a job.
// Returns false if the job is unknown or no snapshot has been recorded.
func (w *Watcher) GetEnv(jobName string) (EnvEntry, bool) {
	w.env.mu.RLock()
	defer w.env.mu.RUnlock()

	entry, ok := w.env.entries[jobName]
	if !ok {
		return EnvEntry{}, false
	}

	copy := make(map[string]string, len(entry.Vars))
	for k, v := range entry.Vars {
		copy[k] = v
	}
	return EnvEntry{Vars: copy}, true
}

// ClearEnv removes the environment snapshot for a job.
func (w *Watcher) ClearEnv(jobName string) error {
	w.env.mu.Lock()
	defer w.env.mu.Unlock()

	if _, ok := w.jobs[jobName]; !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}

	delete(w.env.entries, jobName)
	return nil
}
