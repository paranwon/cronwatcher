package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type nextRunWatcher interface {
	GetNextRun(name string) (time.Time, bool)
	OverdueJobs() []string
}

func makeHandleJobNextRun(w nextRunWatcher) http.HandlerFunc {
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

		next, ok := w.GetNextRun(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}

		payload := map[string]interface{}{
			"job":      name,
			"next_run": next.UTC().Format(time.RFC3339),
			"overdue":  time.Now().After(next),
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(payload)
	}
}
