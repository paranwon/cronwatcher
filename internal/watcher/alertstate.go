package watcher

import "sync"

// AlertState tracks whether an alert has already been fired for a job,
// preventing duplicate notifications on repeated check cycles.
type AlertState struct {
	mu     sync.Mutex
	states map[string]map[string]bool // job -> alertType -> fired
}

func newAlertState(jobs []string) *AlertState {
	s := &AlertState{
		states: make(map[string]map[string]bool, len(jobs)),
	}
	for _, j := range jobs {
		s.states[j] = make(map[string]bool)
	}
	return s
}

// HasFired reports whether the given alertType has already been fired for job.
func (s *AlertState) HasFired(job, alertType string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	types, ok := s.states[job]
	if !ok {
		return false
	}
	return types[alertType]
}

// SetFired marks alertType as fired for job.
func (s *AlertState) SetFired(job, alertType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.states[job]; !ok {
		s.states[job] = make(map[string]bool)
	}
	s.states[job][alertType] = true
}

// Clear resets the alert state for job and alertType, allowing future alerts.
func (s *AlertState) Clear(job, alertType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if types, ok := s.states[job]; ok {
		delete(types, alertType)
	}
}

// ClearAll resets all alert states for job.
func (s *AlertState) ClearAll(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[job] = make(map[string]bool)
}
