package notify

import (
	"context"
	"log"
	"time"
)

// HistoryPruner periodically removes history entries older than RetentionAge.
type HistoryPruner struct {
	watcher prunerWatcher
	interval    time.Duration
	retentionAge time.Duration
	log         *log.Logger
}

type prunerWatcher interface {
	PruneHistory(maxAge time.Duration) int
}

// NewPruner creates a HistoryPruner that runs every interval and removes
// entries older than retentionAge.
func NewPruner(w prunerWatcher, interval, retentionAge time.Duration, logger *log.Logger) *HistoryPruner {
	return &HistoryPruner{
		watcher:     w,
		interval:    interval,
		retentionAge: retentionAge,
		log:         logger,
	}
}

// Run blocks until ctx is cancelled, pruning history on each tick.
func (p *HistoryPruner) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n := p.watcher.PruneHistory(p.retentionAge)
			if n > 0 {
				p.log.Printf("[pruner] removed %d stale history entries", n)
			}
		}
	}
}
