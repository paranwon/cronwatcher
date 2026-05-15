package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type staleWatcher interface {
	MarkStale(now time.Time) ([]watcher.StaleEntry, error)
}

type staleResponse struct {
	JobName    string `json:"job_name"`
	LastSeen   string `json:"last_seen"`
	StaleSince string `json:"stale_since"`
}

func makeHandleJobStale(w staleWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		entries, err := w.MarkStale(time.Now())
		if err != nil {
			http.Error(rw, "internal server error", http.StatusInternalServerError)
			return
		}

		resp := make([]staleResponse, 0, len(entries))
		for _, e := range entries {
			resp = append(resp, staleResponse{
				JobName:    e.JobName,
				LastSeen:   e.LastSeen.UTC().Format(time.RFC3339),
				StaleSince: e.StaleSince.String(),
			})
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(resp)
	}
}
