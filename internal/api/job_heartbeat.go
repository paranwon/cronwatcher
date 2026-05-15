package api

import (
	"net/http"
)

// heartbeatRecorder is the subset of watcher.Watcher used by the heartbeat handler.
type heartbeatRecorder interface {
	RecordHeartbeat(job string) error
}

// makeHandleJobHeartbeat returns an HTTP handler that records a heartbeat for
// the named job. Clients should POST to /jobs/heartbeat?name=<job>.
func makeHandleJobHeartbeat(w heartbeatRecorder) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing required query parameter: name", http.StatusBadRequest)
			return
		}

		if err := w.RecordHeartbeat(name); err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte(`{"status":"ok"}`))
	}
}
