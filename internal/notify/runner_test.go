package notify

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cronwatcher/internal/watcher"
)

func TestRunner_RunsAndStopsOnCancel(t *testing.T) {
	w := watcher.NewForTest()
	w.InjectMissed("cleanup", time.Now().Add(-3*time.Hour))

	m := &mockAlerter{}
	logger := log.New(os.Stdout, "", 0)
	n := New(w, m, logger)

	runner := NewRunner(n, 20*time.Millisecond, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		runner.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// runner exited cleanly
	case <-time.After(200 * time.Millisecond):
		t.Fatal("runner did not stop after context cancellation")
	}

	if len(m.missedCalls) == 0 {
		t.Error("expected at least one missed-job alert to be dispatched")
	}
}

func TestRunner_NoJobs_DoesNotPanic(t *testing.T) {
	w := watcher.NewForTest()
	m := &mockAlerter{}
	logger := log.New(os.Stdout, "", 0)
	n := New(w, m, logger)

	runner := NewRunner(n, 10*time.Millisecond, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	// should not panic
	runner.Run(ctx)
}
