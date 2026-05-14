package notify

import (
	"context"
	"log"
	"time"
)

// Runner periodically triggers missed and long-running job checks.
type Runner struct {
	notifier *Notifier
	interval time.Duration
	logger   *log.Logger
}

// NewRunner creates a Runner that checks on the given interval.
func NewRunner(n *Notifier, interval time.Duration, logger *log.Logger) *Runner {
	return &Runner{
		notifier: n,
		interval: interval,
		logger:   logger,
	}
}

// Run starts the periodic check loop and blocks until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	r.logger.Printf("[notify] runner started, interval=%s", r.interval)

	for {
		select {
		case <-ticker.C:
			r.notifier.CheckMissed()
			r.notifier.CheckLongRunning()
		case <-ctx.Done():
			r.logger.Println("[notify] runner stopped")
			return
		}
	}
}
