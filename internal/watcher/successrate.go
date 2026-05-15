package watcher

import "sync"

// SuccessRate holds computed success rate data for a job.
type SuccessRate struct {
	Total   int     `json:"total"`
	Success int     `json:"success"`
	Failure int     `json:"failure"`
	Rate    float64 `json:"rate"`
}

type successRateStore struct {
	mu   sync.RWMutex
	data map[string]*SuccessRate
}

func newSuccessRateStore(jobs []string) *successRateStore {
	s := &successRateStore{data: make(map[string]*SuccessRate, len(jobs))}
	for _, name := range jobs {
		s.data[name] = &SuccessRate{}
	}
	return s
}

func (s *successRateStore) recordSuccess(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.data[name]
	if !ok {
		return errUnknownJob(name)
	}
	r.Total++
	r.Success++
	r.Rate = float64(r.Success) / float64(r.Total)
	return nil
}

func (s *successRateStore) recordFailure(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.data[name]
	if !ok {
		return errUnknownJob(name)
	}
	r.Total++
	r.Failure++
	r.Rate = float64(r.Success) / float64(r.Total)
	return nil
}

func (s *successRateStore) get(name string) (SuccessRate, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.data[name]
	if !ok {
		return SuccessRate{}, false
	}
	return *r, true
}
