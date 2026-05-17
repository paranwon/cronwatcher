package api

import (
	"encoding/json"
	"net/http"
)

type triggerWatcher interface {
	RecordTrigger(job, triggeredBy, reason string) error
	GetTrigger(job string) (interface{}, bool)
	ClearTrigger(job string) error
}

func makeHandleJobTrigger(w triggerWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			entry, ok := w.GetTrigger(name)
			if !ok {
				http.Error(rw, "not found", http.StatusNotFound)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(rw).Encode(entry)

		case http.MethodPost:
			var body struct {
				TriggeredBy string `json:"triggered_by"`
				Reason      string `json:"reason"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(rw, "invalid body", http.StatusBadRequest)
				return
			}
			if err := w.RecordTrigger(name, body.TriggeredBy, body.Reason); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			if err := w.ClearTrigger(name); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)

		default:
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
