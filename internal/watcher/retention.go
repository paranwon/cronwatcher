package watcher

import "time"

// PruneHistory removes history entries older than maxAge for all known jobs.
// It is safe to call concurrently.
func (w *Watcher) PruneHistory(maxAge time.Duration) int {
	w.mu.Lock()
	defer w.mu.Unlock()

	cutoff := w.clock.Now().Add(-maxAge)
	total := 0

	for name, ring := range w.history {
		filtered := ring[:0]
		for _, entry := range ring {
			if entry.StartedAt.After(cutoff) {
				filtered = append(filtered, entry)
			}
		}
		removed := len(ring) - len(filtered)
		if removed > 0 {
			w.history[name] = filtered
			total += removed
		}
	}

	return total
}
