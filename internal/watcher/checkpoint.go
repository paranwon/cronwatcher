package watcher

import (
	"fmt"
	"sync"
	"time"
)

// CheckpointEntry records a named checkpoint within a job run.
type CheckpointEntry struct {
	Name      string
	RecordedAt time.Time
	Meta      map[string]string
}

type checkpointStore struct {
	mu   sync.RWMutex
	data map[string][]CheckpointEntry // job name -> ordered checkpoints
}

func newCheckpointStore() *checkpointStore {
	return &checkpointStore{
		data: make(map[string][]CheckpointEntry),
	}
}

// RecordCheckpoint appends a named checkpoint for the given job.
func (s *checkpointStore) RecordCheckpoint(job, name string, meta map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[job]; !ok {
		return fmt.Errorf("checkpoint: unknown job %q", job)
	}
	entry := CheckpointEntry{
		Name:       name,
		RecordedAt: time.Now(),
		Meta:       meta,
	}
	s.data[job] = append(s.data[job], entry)
	return nil
}

// GetCheckpoints returns all checkpoints for a job.
func (s *checkpointStore) GetCheckpoints(job string) ([]CheckpointEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, ok := s.data[job]
	if !ok {
		return nil, false
	}
	out := make([]CheckpointEntry, len(entries))
	copy(out, entries)
	return out, true
}

// ClearCheckpoints removes all checkpoints for a job (e.g. on new run start).
func (s *checkpointStore) ClearCheckpoints(job string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[job]; !ok {
		return fmt.Errorf("checkpoint: unknown job %q", job)
	}
	s.data[job] = nil
	return nil
}
