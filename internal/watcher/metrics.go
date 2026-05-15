package watcher

import (
	"sync"
	"time"
)

// JobMetrics holds aggregated runtime statistics for a single job.
type JobMetrics struct {
	JobName     string
	RunCount    int
	ErrorCount  int
	TotalTime   time.Duration
	MinDuration time.Duration
	MaxDuration time.Duration
	LastRun     time.Time
}

// AvgDuration returns the mean run duration, or 0 if no runs have completed.
func (m *JobMetrics) AvgDuration() time.Duration {
	if m.RunCount == 0 {
		return 0
	}
	return m.TotalTime / time.Duration(m.RunCount)
}

type metricsStore struct {
	mu   sync.RWMutex
	data map[string]*JobMetrics
}

func newMetricsStore() *metricsStore {
	return &metricsStore{data: make(map[string]*JobMetrics)}
}

func (s *metricsStore) record(name string, d time.Duration, errored bool, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.data[name]
	if !ok {
		m = &JobMetrics{JobName: name, MinDuration: d, MaxDuration: d}
		s.data[name] = m
	}

	m.RunCount++
	m.TotalTime += d
	m.LastRun = at
	if errored {
		m.ErrorCount++
	}
	if d < m.MinDuration {
		m.MinDuration = d
	}
	if d > m.MaxDuration {
		m.MaxDuration = d
	}
}

func (s *metricsStore) get(name string) (JobMetrics, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.data[name]
	if !ok {
		return JobMetrics{}, false
	}
	copy := *m
	return copy, true
}

func (s *metricsStore) all() []JobMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]JobMetrics, 0, len(s.data))
	for _, m := range s.data {
		copy := *m
		out = append(out, copy)
	}
	return out
}
