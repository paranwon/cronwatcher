package watcher

import (
	"fmt"
	"sync"
	"time"
)

// WindowEntry holds aggregated stats for a rolling time window.
type WindowEntry struct {
	JobName   string
	Window    time.Duration
	Runs      int
	Successes int
	Failures  int
	AvgMs     float64
}

type windowRecord struct {
	ts      time.Time
	success bool
	durMs   float64
}

type windowStore struct {
	mu      sync.Mutex
	records map[string][]windowRecord
}

func newWindowStore(jobs []string) *windowStore {
	m := make(map[string][]windowRecord, len(jobs))
	for _, j := range jobs {
		m[j] = []windowRecord{}
	}
	return &windowStore{records: m}
}

// RecordWindow appends a run result for the given job.
func (s *windowStore) RecordWindow(job string, success bool, dur time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[job]; !ok {
		return fmt.Errorf("window: unknown job %q", job)
	}
	s.records[job] = append(s.records[job], windowRecord{
		ts:      time.Now(),
		success: success,
		durMs:   float64(dur.Milliseconds()),
	})
	return nil
}

// GetWindow returns aggregated stats for runs within the given window duration.
func (s *windowStore) GetWindow(job string, window time.Duration) (WindowEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	recs, ok := s.records[job]
	if !ok {
		return WindowEntry{}, false
	}
	cutoff := time.Now().Add(-window)
	var runs, successes, failures int
	var totalMs float64
	for _, r := range recs {
		if r.ts.Before(cutoff) {
			continue
		}
		runs++
		totalMs += r.durMs
		if r.success {
			successes++
		} else {
			failures++
		}
	}
	var avg float64
	if runs > 0 {
		avg = totalMs / float64(runs)
	}
	return WindowEntry{
		JobName:   job,
		Window:    window,
		Runs:      runs,
		Successes: successes,
		Failures:  failures,
		AvgMs:     avg,
	}, true
}
