package watcher

import (
	"fmt"
	"sync"
	"time"
)

// HeartbeatEntry records the last time a job sent a heartbeat ping.
type HeartbeatEntry struct {
	LastSeen  time.Time
	ExpectsBy time.Time
	Missed    bool
}

type heartbeatStore struct {
	mu      sync.RWMutex
	entries map[string]HeartbeatEntry
}

func newHeartbeatStore() *heartbeatStore {
	return &heartbeatStore{
		entries: make(map[string]HeartbeatEntry),
	}
}

// RecordHeartbeat marks that a job is alive at the current time.
// ttl is how long until the next heartbeat is expected.
func (w *Watcher) RecordHeartbeat(jobName string, ttl time.Duration) error {
	w.mu.RLock()
	_, ok := w.jobs[jobName]
	w.mu.RUnlock()
	if !ok {
		return fmt.Errorf("unknown job: %s", jobName)
	}

	now := time.Now()
	w.heartbeats.mu.Lock()
	w.heartbeats.entries[jobName] = HeartbeatEntry{
		LastSeen:  now,
		ExpectsBy: now.Add(ttl),
		Missed:    false,
	}
	w.heartbeats.mu.Unlock()
	return nil
}

// GetHeartbeat returns the heartbeat entry for a job.
func (w *Watcher) GetHeartbeat(jobName string) (HeartbeatEntry, bool) {
	w.heartbeats.mu.RLock()
	defer w.heartbeats.mu.RUnlock()
	e, ok := w.heartbeats.entries[jobName]
	return e, ok
}

// CheckHeartbeats marks jobs whose heartbeat deadline has passed as missed.
// Returns a list of job names that are newly missed.
func (w *Watcher) CheckHeartbeats() []string {
	now := time.Now()
	w.heartbeats.mu.Lock()
	defer w.heartbeats.mu.Unlock()

	var missed []string
	for name, entry := range w.heartbeats.entries {
		if !entry.Missed && now.After(entry.ExpectsBy) {
			entry.Missed = true
			w.heartbeats.entries[name] = entry
			missed = append(missed, name)
		}
	}
	return missed
}
