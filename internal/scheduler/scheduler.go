package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"cronwatcher/internal/config"
	"cronwatcher/internal/watcher"
)

// Scheduler wraps a cron runner and wires job lifecycle hooks
// into the watcher so missed / long-running jobs are detected.
type Scheduler struct {
	cron    *cron.Cron
	watcher *watcher.Watcher
	cfg     *config.Config
}

// New creates a Scheduler from the provided config and watcher.
func New(cfg *config.Config, w *watcher.Watcher) *Scheduler {
	c := cron.New(cron.WithSeconds())
	return &Scheduler{cron: c, watcher: w, cfg: cfg}
}

// Register adds every job defined in the config to the cron runner.
// Each job is wrapped so the watcher records start/finish times.
func (s *Scheduler) Register() error {
	for _, job := range s.cfg.Jobs {
		job := job // capture loop var
		_, err := s.cron.AddFunc(job.Schedule, func() {
			start := time.Now()
			s.watcher.RecordStart(job.Name, start)
			log.Printf("[scheduler] job started: %s", job.Name)

			// Placeholder: real implementations would exec a command here.
			// The watcher's ticker handles long-running detection independently.

			s.watcher.RecordFinish(job.Name, time.Now())
			log.Printf("[scheduler] job finished: %s (duration: %s)",
				job.Name, time.Since(start))
		})
		if err != nil {
			return err
		}
		log.Printf("[scheduler] registered job %q with schedule %q", job.Name, job.Schedule)
	}
	return nil
}

// Start begins the cron scheduler.
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("[scheduler] started")
}

// Stop gracefully stops the cron scheduler.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("[scheduler] stopped")
}
