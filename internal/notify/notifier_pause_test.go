package notify_test

import (
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/config"
	"github.com/cronwatcher/cronwatcher/internal/notify"
	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestCheckMissed_SkipsPausedJob(t *testing.T) {
	cfg := &config.Config{
		Jobs: []config.Job{
			{Name: "paused-job", Schedule: "* * * * *", MaxDuration: "5m", GracePeriod: "1m"},
		},
	}

	w := watcher.NewForTest(cfg)
	// Mark the job as paused before the notifier checks it.
	_ = w.PauseJob("paused-job")

	var alerted bool
	n := newTestNotifier(cfg, w, func(name string) {
		alerted = true
		_ = name
	})

	// Simulate the job being overdue.
	w.RecordStart("paused-job")
	time.Sleep(5 * time.Millisecond)

	n.CheckMissed(time.Now().Add(-time.Hour))

	if alerted {
		t.Error("expected no alert for paused job, but got one")
	}
}
