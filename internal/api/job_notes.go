package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type notesWatcher interface {
	AddNote(job, text string) (watcher.NoteEntry, error)
	GetNotes(job string) ([]watcher.NoteEntry, bool)
	DeleteNote(job, id string) error
}

func makeHandleJobNotes(w notesWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			entries, ok := w.GetNotes(name)
			if !ok {
				http.Error(rw, "job not found", http.StatusNotFound)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(entries)

		case http.MethodPost:
			var body struct {
				Text string `json:"text"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Text == "" {
				http.Error(rw, "invalid body", http.StatusBadRequest)
				return
			}
			entry, err := w.AddNote(name, body.Text)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusCreated)
			json.NewEncoder(rw).Encode(entry)

		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(rw, "missing id", http.StatusBadRequest)
				return
			}
			if err := w.DeleteNote(name, id); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)

		default:
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
