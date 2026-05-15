package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/api"
)

func TestHandleJobPause_PausesJob(t *testing.T) {
	w := newTestWatcher()
	s := api.New(":0", w, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/pause?name=test-job&action=pause", nil)
	rec := httptest.NewRecorder()
	s.HandleJobPause(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body["status"] != "paused" {
		t.Errorf("expected status=paused, got %q", body["status"])
	}
}

func TestHandleJobPause_ResumesJob(t *testing.T) {
	w := newTestWatcher()
	s := api.New(":0", w, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/pause?name=test-job&action=resume", nil)
	rec := httptest.NewRecorder()
	s.HandleJobPause(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandleJobPause_MissingName_ReturnsBadRequest(t *testing.T) {
	w := newTestWatcher()
	s := api.New(":0", w, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/pause?action=pause", nil)
	rec := httptest.NewRecorder()
	s.HandleJobPause(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandleJobPause_InvalidAction_ReturnsBadRequest(t *testing.T) {
	w := newTestWatcher()
	s := api.New(":0", w, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/pause?name=test-job&action=stop", nil)
	rec := httptest.NewRecorder()
	s.HandleJobPause(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandleJobPause_MethodNotAllowed(t *testing.T) {
	w := newTestWatcher()
	s := api.New(":0", w, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs/pause?name=test-job&action=pause", nil)
	rec := httptest.NewRecorder()
	s.HandleJobPause(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
