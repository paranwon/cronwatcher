package notify_test

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dop251/cronwatcher/internal/notify"
)

type fakePrunerWatcher struct {
	calls atomic.Int32
	returns int
}

func (f *fakePrunerWatcher) PruneHistory(_ time.Duration) int {
	f.calls.Add(1)
	return f.returns
}

func TestPruner_CallsPruneOnTick(t *testing.T) {
	fw := &fakePrunerWatcher{returns: 3}
	logger := log.New(os.Stderr, "", 0)

	p := notify.NewPruner(fw, 20*time.Millisecond, 24*time.Hour, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()

	p.Run(ctx)

	if fw.calls.Load() < 2 {
		t.Fatalf("expected at least 2 prune calls, got %d", fw.calls.Load())
	}
}

func TestPruner_StopsOnContextCancel(t *testing.T) {
	fw := &fakePrunerWatcher{returns: 0}
	logger := log.New(os.Stderr, "", 0)

	p := notify.NewPruner(fw, 1*time.Second, 24*time.Hour, logger)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		p.Run(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("pruner did not stop after context cancel")
	}
}
