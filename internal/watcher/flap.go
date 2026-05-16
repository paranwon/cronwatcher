package watcher

import (
	"fmt"
	"sync"
	"time"
)

// FlapEntry records rapid state changes (success/failure alternations) for a job.
type FlapEntry struct {
	JobName    string
	Flaps      int
	WindowStart time.Time
	LastSeen   time.Time
}

type flapStore struct {
	mu      sync.Mutex
	entries map[string]*flapEntry
}

type flapEntry struct {
	count       int
	windowStart time.Time
	lastState   string // "success" or "failure"
	lastSeen    time.Time
}

func newFlapStore() *flapStore {
	return &flapStore{entries: make(map[string]*flapEntry)}
}

// RecordFlapState records a state transition and increments the flap counter if
// the state alternates within the detection window.
func (f *flapStore) RecordFlapState(job, state string, window time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	e, ok := f.entries[job]
	if !ok {
		return fmt.Errorf("flap: unknown job %q", job)
	}

	now := time.Now()

	// Reset window if expired.
	if now.Sub(e.windowStart) > window {
		e.count = 0
		e.windowStart = now
		e.lastState = ""
	}

	if e.lastState != "" && e.lastState != state {
		e.count++
	}

	e.lastState = state
	e.lastSeen = now
	return nil
}

// GetFlap returns the current flap entry for a job.
func (f *flapStore) GetFlap(job string) (FlapEntry, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	e, ok := f.entries[job]
	if !ok {
		return FlapEntry{}, false
	}
	return FlapEntry{
		JobName:     job,
		Flaps:       e.count,
		WindowStart: e.windowStart,
		LastSeen:    e.lastSeen,
	}, true
}

// registerFlapJob registers a job in the flap store.
func (f *flapStore) registerFlapJob(job string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entries[job] = &flapEntry{windowStart: time.Now()}
}
