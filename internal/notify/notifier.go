package notify

import (
	"fmt"
	"log"
	"time"

	"github.com/cronwatcher/internal/alert"
	"github.com/cronwatcher/internal/watcher"
)

// Notifier watches for missed and long-running jobs and dispatches alerts.
type Notifier struct {
	w       *watcher.Watcher
	alerter alert.Alerter
	logger  *log.Logger
}

// New creates a new Notifier.
func New(w *watcher.Watcher, a alert.Alerter, logger *log.Logger) *Notifier {
	return &Notifier{
		w:       w,
		alerter: a,
		logger:  logger,
	}
}

// CheckMissed evaluates all job statuses and fires alerts for any missed jobs.
func (n *Notifier) CheckMissed() {
	statuses := n.w.Statuses()
	for _, s := range statuses {
		if s.Missed {
			n.logger.Printf("[notify] missed job detected: %s", s.Name)
			if err := n.alerter.MissedJob(s.Name, s.LastRun); err != nil {
				n.logger.Printf("[notify] alert error for missed job %s: %v", s.Name, err)
			}
		}
	}
}

// CheckLongRunning evaluates all job statuses and fires alerts for long-running jobs.
func (n *Notifier) CheckLongRunning() {
	statuses := n.w.Statuses()
	for _, s := range statuses {
		if s.Running && s.Duration > 0 && s.ExceedsThreshold {
			n.logger.Printf("[notify] long-running job detected: %s (duration: %s)", s.Name, formatDuration(s.Duration))
			if err := n.alerter.LongRunningJob(s.Name, s.Duration); err != nil {
				n.logger.Printf("[notify] alert error for long-running job %s: %v", s.Name, err)
			}
		}
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
}
