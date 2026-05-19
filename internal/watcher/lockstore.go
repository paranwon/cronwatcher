package watcher

import (
	"fmt"
	"sync"
	"time"
)

// LockEntry represents an acquired distributed-style lock for a job.
type LockEntry struct {
	Owner     string    `json:"owner"`
	AcquiredAt time.Time `json:"acquired_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// IsExpired reports whether the lock has passed its expiry time.
func (l LockEntry) IsExpired(now time.Time) bool {
	return now.After(l.ExpiresAt)
}

type lockStore struct {
	mu    sync.Mutex
	locks map[string]LockEntry // keyed by job name
}

func newLockStore() *lockStore {
	return &lockStore{locks: make(map[string]LockEntry)}
}

// AcquireLock attempts to acquire a lock for the given job.
// Returns an error if the job is unknown or a non-expired lock already exists.
func (s *lockStore) AcquireLock(job, owner string, ttl time.Duration, known func(string) bool, now time.Time) error {
	if !known(job) {
		return fmt.Errorf("lockstore: unknown job %q", job)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.locks[job]; ok && !existing.IsExpired(now) {
		return fmt.Errorf("lockstore: job %q is already locked by %q", job, existing.Owner)
	}
	s.locks[job] = LockEntry{
		Owner:      owner,
		AcquiredAt: now,
		ExpiresAt:  now.Add(ttl),
	}
	return nil
}

// ReleaseLock releases the lock for a job. Returns an error if no lock is held.
func (s *lockStore) ReleaseLock(job string, known func(string) bool) error {
	if !known(job) {
		return fmt.Errorf("lockstore: unknown job %q", job)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.locks[job]; !ok {
		return fmt.Errorf("lockstore: no lock held for job %q", job)
	}
	delete(s.locks, job)
	return nil
}

// GetLock returns the current LockEntry for a job, if one exists.
func (s *lockStore) GetLock(job string) (LockEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.locks[job]
	return e, ok
}
