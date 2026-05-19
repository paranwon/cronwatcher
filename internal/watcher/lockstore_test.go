package watcher

import (
	"testing"
	"time"
)

func knownJobs(names ...string) func(string) bool {
	set := make(map[string]struct{}, len(names))
	for _, n := range names {
		set[n] = struct{}{}
	}
	return func(job string) bool {
		_, ok := set[job]
		return ok
	}
}

func TestAcquireLock_StoresEntry(t *testing.T) {
	s := newLockStore()
	now := time.Now()
	known := knownJobs("backup")

	if err := s.AcquireLock("backup", "worker-1", 30*time.Second, known, now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry, ok := s.GetLock("backup")
	if !ok {
		t.Fatal("expected lock entry to exist")
	}
	if entry.Owner != "worker-1" {
		t.Errorf("expected owner worker-1, got %q", entry.Owner)
	}
	if !entry.ExpiresAt.Equal(now.Add(30 * time.Second)) {
		t.Errorf("unexpected expires_at: %v", entry.ExpiresAt)
	}
}

func TestAcquireLock_UnknownJob_ReturnsError(t *testing.T) {
	s := newLockStore()
	err := s.AcquireLock("ghost", "worker-1", time.Minute, knownJobs(), time.Now())
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestAcquireLock_AlreadyLocked_ReturnsError(t *testing.T) {
	s := newLockStore()
	now := time.Now()
	known := knownJobs("sync")

	_ = s.AcquireLock("sync", "worker-1", time.Minute, known, now)
	err := s.AcquireLock("sync", "worker-2", time.Minute, known, now)
	if err == nil {
		t.Fatal("expected error when lock already held")
	}
}

func TestAcquireLock_ExpiredLock_AllowsReacquire(t *testing.T) {
	s := newLockStore()
	known := knownJobs("report")
	past := time.Now().Add(-2 * time.Minute)

	_ = s.AcquireLock("report", "worker-old", time.Second, known, past)

	// lock has expired; new acquire should succeed
	if err := s.AcquireLock("report", "worker-new", time.Minute, known, time.Now()); err != nil {
		t.Fatalf("expected reacquire to succeed, got: %v", err)
	}
	entry, _ := s.GetLock("report")
	if entry.Owner != "worker-new" {
		t.Errorf("expected owner worker-new, got %q", entry.Owner)
	}
}

func TestReleaseLock_RemovesEntry(t *testing.T) {
	s := newLockStore()
	known := knownJobs("cleanup")
	_ = s.AcquireLock("cleanup", "worker-1", time.Minute, known, time.Now())

	if err := s.ReleaseLock("cleanup", known); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.GetLock("cleanup")
	if ok {
		t.Error("expected lock to be removed")
	}
}

func TestReleaseLock_NoLock_ReturnsError(t *testing.T) {
	s := newLockStore()
	known := knownJobs("idle")
	err := s.ReleaseLock("idle", known)
	if err == nil {
		t.Fatal("expected error when no lock held")
	}
}

func TestGetLock_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newLockStore()
	_, ok := s.GetLock("nonexistent")
	if ok {
		t.Error("expected false for unknown job")
	}
}
