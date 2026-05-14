package watcher

import (
	"fmt"
	"time"

	"github.com/cronwatcher/internal/api"
)

// Status returns a snapshot of all tracked job statuses.
func (w *Watcher) Status() map[string]api.JobStatus {
	w.mu.RLock()
	defer w.mu.RUnlock()

	result := make(map[string]api.JobStatus, len(w.jobs))
	for name, state := range w.jobs {
		js := api.JobStatus{
			Name:    name,
			Running: state.running,
		}
		if !state.lastStart.IsZero() {
			t := state.lastStart
			js.LastStart = &t
		}
		if !state.lastEnd.IsZero() {
			t := state.lastEnd
			js.LastEnd = &t
			dur := formatDuration(state.lastEnd.Sub(state.lastStart))
			js.LastDuration = &dur
		}
		result[name] = js
	}
	return result
}

// ServeHTTP allows *Server to satisfy http.Handler for testing without a real listener.
func (s *api.Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.Handler.ServeHTTP(w, r)
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%.2fs", int(d.Minutes()), d.Seconds()-float64(int(d.Minutes()))*60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
