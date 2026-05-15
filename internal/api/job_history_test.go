package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/densestvoid/cronwatcher/internal/watcher"
)

func TestHandleJobHistory_ReturnsHistory(t *testing.T) {
	w := newTestWatcher()
	w.RecordStart("backup")
	time.Sleep(2 * time.Millisecond)
	w.RecordFinish("backup", nil)

	srv := New(w, nil)
	req := httptest.NewRequest(http.MethodGet, "/job/history?name=backup", nil)
	rec := httptest.NewRecorder()
	srv.handleJobHistory(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var entries []watcher.HistoryEntry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one history entry")
	}
}

func TestHandleJobHistory_NotFound(t *testing.T) {
	w := newTestWatcher()
	srv := New(w, nil)
	req := httptest.NewRequest(http.MethodGet, "/job/history?name=unknown", nil)
	rec := httptest.NewRecorder()
	srv.handleJobHistory(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleJobHistory_MissingName(t *testing.T) {
	w := newTestWatcher()
	srv := New(w, nil)
	req := httptest.NewRequest(http.MethodGet, "/job/history", nil)
	rec := httptest.NewRecorder()
	srv.handleJobHistory(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleJobHistory_InvalidLimit(t *testing.T) {
	w := newTestWatcher()
	srv := New(w, nil)
	req := httptest.NewRequest(http.MethodGet, "/job/history?name=backup&limit=abc", nil)
	rec := httptest.NewRecorder()
	srv.handleJobHistory(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleJobHistory_MethodNotAllowed(t *testing.T) {
	w := newTestWatcher()
	srv := New(w, nil)
	req := httptest.NewRequest(http.MethodPost, "/job/history?name=backup", nil)
	rec := httptest.NewRecorder()
	srv.handleJobHistory(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
