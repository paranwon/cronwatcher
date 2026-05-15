package watcher

import "sync"

// StreakEntry holds consecutive success/failure counts for a job.
type StreakEntry struct {
	SuccessStreak int
	FailureStreak int
}

type streakStore struct {
	mu      sync.RWMutex
	streaks map[string]*StreakEntry
}

func newStreakStore(jobs []string) *streakStore {
	s := &streakStore{streaks: make(map[string]*StreakEntry, len(jobs))}
	for _, j := range jobs {
		s.streaks[j] = &StreakEntry{}
	}
	return s
}

// RecordStreakSuccess increments the success streak and resets the failure streak.
func (w *Watcher) RecordStreakSuccess(name string) error {
	w.streaks.mu.Lock()
	defer w.streaks.mu.Unlock()
	e, ok := w.streaks.streaks[name]
	if !ok {
		return ErrUnknownJob
	}
	e.SuccessStreak++
	e.FailureStreak = 0
	return nil
}

// RecordStreakFailure increments the failure streak and resets the success streak.
func (w *Watcher) RecordStreakFailure(name string) error {
	w.streaks.mu.Lock()
	defer w.streaks.mu.Unlock()
	e, ok := w.streaks.streaks[name]
	if !ok {
		return ErrUnknownJob
	}
	e.FailureStreak++
	e.SuccessStreak = 0
	return nil
}

// GetStreak returns the streak entry for a job.
func (w *Watcher) GetStreak(name string) (StreakEntry, bool) {
	w.streaks.mu.RLock()
	defer w.streaks.mu.RUnlock()
	e, ok := w.streaks.streaks[name]
	if !ok {
		return StreakEntry{}, false
	}
	return *e, true
}
