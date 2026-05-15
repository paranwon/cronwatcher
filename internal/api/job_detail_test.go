package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnnrly/cronwatcher/internal/watcher"
)

func TestHandleJobDetail_ReturnsJob(t *testing.T) {
	w := newTestWatcher()
	w.RecordStart("backup")

	s := New(w, nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/backup", nil)
	rec := httptest.NewRecorder()

	s.handleJobDetail(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp jobDetailResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Name != "backup" {
		t.Errorf("expected name 'backup', got %q", resp.Name)
	}
	if !resp.Running {
		t.Errorf("expected running=true")
	}
}

func TestHandleJobDetail_NotFound(t *testing.T) {
	w := newTestWatcher()
	s := New(w, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs/nonexistent", nil)
	rec := httptest.NewRecorder()

	s.handleJobDetail(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleJobDetail_MissingName(t *testing.T) {
	w := newTestWatcher()
	s := New(w, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs/", nil)
	rec := httptest.NewRecorder()

	s.handleJobDetail(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleJobDetail_MethodNotAllowed(t *testing.T) {
	w := newTestWatcher()
	s := New(w, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/backup", nil)
	rec := httptest.NewRecorder()

	s.handleJobDetail(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandleJobDetail_ShowsLastSeen(t *testing.T) {
	w := newTestWatcher()
	w.RecordStart("cleanup")
	time.Sleep(2 * time.Millisecond)
	w.RecordFinish("cleanup")

	s := New(w, nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/cleanup", nil)
	rec := httptest.NewRecorder()

	s.handleJobDetail(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp jobDetailResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Running {
		t.Errorf("expected running=false after finish")
	}
	if resp.LastSeen == "" {
		t.Errorf("expected last_seen to be set")
	}
	if resp.LastElapsed == "" {
		t.Errorf("expected last_elapsed to be set")
	}
}

// Ensure watcher.Watcher satisfies the interface used by newTestWatcher.
var _ interface {
	Status() []watcher.JobStatus
	RecordStart(string)
	RecordFinish(string)
} = (*watcher.Watcher)(nil)
