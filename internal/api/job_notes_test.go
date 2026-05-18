package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type fakeNotesWatcher struct {
	notes map[string][]watcher.NoteEntry
}

func newFakeNotesWatcher() *fakeNotesWatcher {
	return &fakeNotesWatcher{notes: map[string][]watcher.NoteEntry{"job1": {}}}
}

func (f *fakeNotesWatcher) AddNote(job, text string) (watcher.NoteEntry, error) {
	w := watcher.NewForTest([]string{job})
	return w.AddNote(job, text)
}

func (f *fakeNotesWatcher) GetNotes(job string) ([]watcher.NoteEntry, bool) {
	entries, ok := f.notes[job]
	return entries, ok
}

func (f *fakeNotesWatcher) DeleteNote(job, id string) error {
	w := watcher.NewForTest([]string{job})
	_, _ = w.AddNote(job, "placeholder")
	return w.DeleteNote(job, id)
}

func TestHandleJobNotes_GET_ReturnsNotes(t *testing.T) {
	fw := newFakeNotesWatcher()
	fw.notes["job1"] = []watcher.NoteEntry{{ID: "1", Text: "hello"}}
	h := makeHandleJobNotes(fw)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notes?name=job1", nil)
	h(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []watcher.NoteEntry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 1 || entries[0].Text != "hello" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}

func TestHandleJobNotes_GET_NotFound(t *testing.T) {
	h := makeHandleJobNotes(newFakeNotesWatcher())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notes?name=ghost", nil)
	h(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestHandleJobNotes_POST_CreatesNote(t *testing.T) {
	w := watcher.NewForTest([]string{"job1"})
	h := makeHandleJobNotes(w)
	body, _ := json.Marshal(map[string]string{"text": "check logs"})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notes?name=job1", bytes.NewReader(body))
	h(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var entry watcher.NoteEntry
	if err := json.NewDecoder(rr.Body).Decode(&entry); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if entry.Text != "check logs" {
		t.Errorf("unexpected text: %s", entry.Text)
	}
}

func TestHandleJobNotes_MissingName_ReturnsBadRequest(t *testing.T) {
	h := makeHandleJobNotes(newFakeNotesWatcher())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	h(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandleJobNotes_MethodNotAllowed(t *testing.T) {
	h := makeHandleJobNotes(newFakeNotesWatcher())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/notes?name=job1", nil)
	h(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
