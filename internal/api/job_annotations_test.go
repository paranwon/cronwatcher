package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/api"
	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func newAnnotationWatcher() *watcher.Watcher {
	return watcher.NewForTest([]string{"backup"})
}

func TestHandleJobAnnotations_GET_ReturnsAnnotations(t *testing.T) {
	w := newAnnotationWatcher()
	_ = w.SetAnnotations("backup", map[string]string{"owner": "ops"})
	srv := api.New(w, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/annotations?name=backup", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var anns map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&anns)
	if anns["owner"] != "ops" {
		t.Errorf("unexpected annotations: %v", anns)
	}
}

func TestHandleJobAnnotations_GET_NotFound(t *testing.T) {
	w := newAnnotationWatcher()
	srv := api.New(w, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/annotations?name=ghost", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleJobAnnotations_POST_SetsAnnotations(t *testing.T) {
	w := newAnnotationWatcher()
	srv := api.New(w, nil)

	body, _ := json.Marshal(map[string]string{"env": "prod"})
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/annotations?name=backup", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	anns, _ := w.GetAnnotations("backup")
	if anns["env"] != "prod" {
		t.Errorf("expected annotation to be set, got: %v", anns)
	}
}

func TestHandleJobAnnotations_DELETE_RemovesKey(t *testing.T) {
	w := newAnnotationWatcher()
	_ = w.SetAnnotations("backup", map[string]string{"owner": "ops"})
	srv := api.New(w, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/jobs/annotations?name=backup&key=owner", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	anns, _ := w.GetAnnotations("backup")
	if _, exists := anns["owner"]; exists {
		t.Error("expected annotation to be deleted")
	}
}

func TestHandleJobAnnotations_MissingName_ReturnsBadRequest(t *testing.T) {
	w := newAnnotationWatcher()
	srv := api.New(w, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/annotations", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleJobAnnotations_MethodNotAllowed(t *testing.T) {
	w := newAnnotationWatcher()
	srv := api.New(w, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/jobs/annotations?name=backup", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
