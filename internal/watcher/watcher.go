package watcher

import (
	"sync"
	"time"

	"github.com/user/cronwatcher/internal/config"
)

// JobState tracks the runtime state of a monitored cron job.
type JobState struct {
	Name      string
	LastSeen  time.Time
	StartedAt time.Time
	Running   bool
	Mu        sync.Mutex
}

// Watcher monitors cron job heartbeats and detects missed or long-running jobs.
type Watcher struct {
	cfg    *config.Config
	states map[string]*JobState
	alerts AlertSink
	stopCh chan struct{}
	mu     sync.RWMutex
}

// AlertSink is implemented by anything that can receive watcher alerts.
type AlertSink interface {
	MissedJob(job config.Job, lastSeen time.Time)
	LongRunningJob(job config.Job, duration time.Duration)
}

// New creates a Watcher with the given config and alert sink.
func New(cfg *config.Config, alerts AlertSink) *Watcher {
	states := make(map[string]*JobState, len(cfg.Jobs))
	for _, j := range cfg.Jobs {
		states[j.Name] = &JobState{Name: j.Name}
	}
	return &Watcher{
		cfg:    cfg,
		states: states,
		alerts: alerts,
		stopCh: make(chan struct{}),
	}
}

// Start begins the background polling loop.
func (w *Watcher) Start() {
	ticker := time.NewTicker(w.cfg.CheckInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.check(time.Now())
			case <-w.stopCh:
				return
			}
		}
	}()
}

// Stop halts the background polling loop.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

// RecordStart marks a job as started at the given time.
func (w *Watcher) RecordStart(name string, at time.Time) {
	w.mu.RLock()
	s, ok := w.states[name]
	w.mu.RUnlock()
	if !ok {
		return
	}
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Running = true
	s.StartedAt = at
	s.LastSeen = at
}

// RecordFinish marks a job as finished at the given time.
func (w *Watcher) RecordFinish(name string, at time.Time) {
	w.mu.RLock()
	s, ok := w.states[name]
	w.mu.RUnlock()
	if !ok {
		return
	}
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Running = false
	s.LastSeen = at
}

// check evaluates all jobs against their configured thresholds.
func (w *Watcher) check(now time.Time) {
	for _, job := range w.cfg.Jobs {
		w.mu.RLock()
		s := w.states[job.Name]
		w.mu.RUnlock()

		s.Mu.Lock()
		lastSeen := s.LastSeen
		running := s.Running
		startedAt := s.StartedAt
		s.Mu.Unlock()

		if !lastSeen.IsZero() && !running {
			if now.Sub(lastSeen) > job.MaxDelay {
				w.alerts.MissedJob(job, lastSeen)
			}
		}
		if running {
			if now.Sub(startedAt) > job.MaxDuration {
				w.alerts.LongRunningJob(job, now.Sub(startedAt))
			}
		}
	}
}
