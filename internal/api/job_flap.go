package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// FlapReader is the interface required to read flap state for a job.
type FlapReader interface {
	GetFlap(job string) (FlapEntry, bool)
}

// FlapEntry mirrors watcher.FlapEntry for the API response.
type FlapEntry struct {
	JobName     string    `json:"job_name"`
	Flaps       int       `json:"flaps"`
	WindowStart time.Time `json:"window_start"`
	LastSeen    time.Time `json:"last_seen"`
}

func makeHandleJobFlap(w FlapReader) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		entry, ok := w.GetFlap(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(FlapEntry{
			JobName:     entry.JobName,
			Flaps:       entry.Flaps,
			WindowStart: entry.WindowStart,
			LastSeen:    entry.LastSeen,
		})
	}
}
