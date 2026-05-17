package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type exitCodeWatcher interface {
	RecordExitCode(job string, code int) error
	GetExitCode(job string) (exitCodeEntry, bool)
}

type exitCodeEntry struct {
	Code       int
	RecordedAt time.Time
}

func makeHandleJobExitCode(w exitCodeWatcher) http.HandlerFunc {
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
				Code int `json:"code"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(rw, "invalid body", http.StatusBadRequest)
				return
			}
			if err := w.RecordExitCode(name, body.Code); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		e, ok := w.GetExitCode(name)
		if !ok {
			http.Error(rw, "not found", http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"code":        e.Code,
			"recorded_at": e.RecordedAt,
		})
	}
}
