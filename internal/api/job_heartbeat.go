package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type heartbeatWatcher interface {
	RecordHeartbeat(jobName string, ttl time.Duration) error
}

// makeHandleJobHeartbeat returns an HTTP handler that accepts a heartbeat ping
// for a named job. Query params: name (required), ttl (seconds, default 300).
func makeHandleJobHeartbeat(w heartbeatWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing required query param: name", http.StatusBadRequest)
			return
		}

		ttlSecs := 300
		if raw := r.URL.Query().Get("ttl"); raw != "" {
			v, err := strconv.Atoi(raw)
			if err != nil || v <= 0 {
				http.Error(rw, "ttl must be a positive integer (seconds)", http.StatusBadRequest)
				return
			}
			ttlSecs = v
		}

		ttl := time.Duration(ttlSecs) * time.Second
		if err := w.RecordHeartbeat(name, ttl); err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(rw).Encode(map[string]interface{}{
			"job":        name,
			"ttl_seconds": ttlSecs,
			"recorded":   true,
		})
	}
}
