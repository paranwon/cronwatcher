package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type slaWatcher interface {
	RecordSLA(name string, deadline time.Time, met bool) error
	GetSLA(name string) (interface{ GetMetSLA() bool }, bool)
}

type slaGetter interface {
	GetSLA(name string) (slaEntry, bool)
}

type slaEntry struct {
	JobName   string    `json:"job_name"`
	Deadline  time.Time `json:"deadline"`
	MetSLA    bool      `json:"met_sla"`
	CheckedAt time.Time `json:"checked_at"`
}

type slaStatusWatcher interface {
	GetSLA(name string) (slaEntry, bool)
}

func makeHandleJobSLA(w slaStatusWatcher) http.HandlerFunc {
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
		entry, ok := w.GetSLA(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(entry)
	}
}
