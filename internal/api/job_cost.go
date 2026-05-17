package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/densestvoid/cronwatcher/internal/watcher"
)

type costWatcher interface {
	RecordCost(job string, e watcher.CostEntry) error
	GetCostSummary(job string) (watcher.CostSummary, bool)
}

func makeHandleJobCost(w costWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			var body struct {
				CPUSecs  float64 `json:"cpu_secs"`
				MemoryMB float64 `json:"memory_mb"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(rw, "invalid body", http.StatusBadRequest)
				return
			}
			err := w.RecordCost(name, watcher.CostEntry{
				RecordedAt: time.Now(),
				CPUSecs:    body.CPUSecs,
				MemoryMB:   body.MemoryMB,
			})
			if err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		sum, ok := w.GetCostSummary(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(sum)
	}
}
