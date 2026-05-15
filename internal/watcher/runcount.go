package watcher

import "sync"

// runCountStore tracks the number of successful and failed runs per job.
type runCountStore struct {
	mu    sync.RWMutex
	counts map[string]*RunCount
}

// RunCount holds success and failure totals for a job.
type RunCount struct {
	Success int64
	Failure int64
	Total   int64
}

func newRunCountStore() *runCountStore {
	return &runCountStore{
		counts: make(map[string]*RunCount),
	}
}

func (s *runCountStore) recordSuccess(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.counts[name]
	if !ok {
		return errJobNotFound(name)
	}
	c.Success++
	c.Total++
	return nil
}

func (s *runCountStore) recordFailure(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.counts[name]
	if !ok {
		return errJobNotFound(name)
	}
	c.Failure++
	c.Total++
	return nil
}

func (s *runCountStore) getRunCount(name string) (RunCount, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.counts[name]
	if !ok {
		return RunCount{}, false
	}
	return *c, true
}

func (s *runCountStore) register(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.counts[name]; !ok {
		s.counts[name] = &RunCount{}
	}
}
