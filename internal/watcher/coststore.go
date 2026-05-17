package watcher

import (
	"fmt"
	"sync"
	"time"
)

// CostEntry records the resource cost of a single job run.
type CostEntry struct {
	RecordedAt time.Time
	Duration   time.Duration
	CPUSecs    float64
	MemoryMB   float64
}

// CostSummary aggregates cost data across all recorded runs.
type CostSummary struct {
	RunCount   int
	TotalCPU   float64
	TotalMemMB float64
	AvgCPU     float64
	AvgMemMB   float64
}

type costStore struct {
	mu      sync.RWMutex
	entries map[string][]CostEntry
}

func newCostStore(jobs []string) *costStore {
	m := make(map[string][]CostEntry, len(jobs))
	for _, j := range jobs {
		m[j] = nil
	}
	return &costStore{entries: m}
}

func (s *costStore) recordCost(job string, e CostEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[job]; !ok {
		return fmt.Errorf("costStore: unknown job %q", job)
	}
	s.entries[job] = append(s.entries[job], e)
	return nil
}

func (s *costStore) getCostSummary(job string) (CostSummary, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, ok := s.entries[job]
	if !ok {
		return CostSummary{}, false
	}
	if len(entries) == 0 {
		return CostSummary{}, true
	}
	var totalCPU, totalMem float64
	for _, e := range entries {
		totalCPU += e.CPUSecs
		totalMem += e.MemoryMB
	}
	n := float64(len(entries))
	return CostSummary{
		RunCount:   len(entries),
		TotalCPU:   totalCPU,
		TotalMemMB: totalMem,
		AvgCPU:     totalCPU / n,
		AvgMemMB:   totalMem / n,
	}, true
}
