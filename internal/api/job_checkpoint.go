package api

import (
	"encoding/json"
	"net/http"

	"github.com/densestvoid/cronwatcher/internal/watcher"
)

type checkpointWatcher interface {
	RecordCheckpoint(job, name string, meta map[string]string) error
	GetCheckpoints(job string) ([]watcher.CheckpointEntry, bool)
	ClearCheckpoints(job string) error
}

func makeHandleJobCheckpoint(w checkpointWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodDelete {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			entries, ok := w.GetCheckpoints(name)
			if !ok {
				http.Error(rw, "job not found", http.StatusNotFound)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(entries)
		case http.MethodPost:
			var body struct {
				Checkpoint string            `json:"checkpoint"`
				Meta       map[string]string `json:"meta"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Checkpoint == "" {
				http.Error(rw, "invalid body", http.StatusBadRequest)
				return
			}
			if err := w.RecordCheckpoint(name, body.Checkpoint, body.Meta); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			if err := w.ClearCheckpoints(name); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
		}
	}
}
