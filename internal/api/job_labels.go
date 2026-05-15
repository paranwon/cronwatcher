package api

import (
	"encoding/json"
	"net/http"
)

type labelWatcher interface {
	SetLabels(jobName string, labels map[string]string) error
	GetLabels(jobName string) (map[string]string, bool)
	DeleteLabel(jobName, key string) error
}

func makeHandleJobLabels(w labelWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			labels, ok := w.GetLabels(name)
			if !ok {
				http.Error(rw, "job not found", http.StatusNotFound)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(labels)

		case http.MethodPost:
			var labels map[string]string
			if err := json.NewDecoder(r.Body).Decode(&labels); err != nil {
				http.Error(rw, "invalid body", http.StatusBadRequest)
				return
			}
			if err := w.SetLabels(name, labels); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(rw, "missing key", http.StatusBadRequest)
				return
			}
			if err := w.DeleteLabel(name, key); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)

		default:
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
