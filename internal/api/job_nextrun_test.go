package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeNextRunWatcher struct {
	entries map[string]time.Time
}

func (f *fakeNextRunWatcher) GetNextRun(name string) (time.Time, bool) {
	t, ok := f.entries[name]
	return t, ok
}

func (f *fakeNextRunWatcher) OverdueJobs() []string {
	var out []string
	for name, t := range f.entries {
		if time.Now().After(t) {
			out = append(out, name)
		}
	}
	return out
}

func newFakeNextRunWatcher() *fakeNextRunWatcher {
	return &fakeNextRunWatcher{entries: make(map[string]time.Time)}
}

func TestHandleJobNextRun_ReturnsNextRun(t *testing.T) {
	fw := newFakeNextRunWatcher()
	future := time.Now().Add(10 * time.Minute)
	fw.entries["backup"] = future

	h := makeHandleJobNextRun(fw)
	req := httptest.NewRequest(http.MethodGet, "/jobs/nextrun?name=backup", nil)
	rr := httptest.NewRecorder()
	h(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["job"] != "backup" {
		t.Errorf("expected job=backup, got %v", body["job"])
	}
	if body["overdue"] != false {
		t.Errorf("expected overdue=false, got %v", body["overdue"])
	}
}

func TestHandleJobNextRun_MissingName_ReturnsBadRequest(t *testing.T) {
	h := makeHandleJobNextRun(newFakeNextRunWatcher())
	req := httptest.NewRequest(http.MethodGet, "/jobs/nextrun", nil)
	rr := httptest.NewRecorder()
	h(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHandleJobNextRun_NotFound_Returns404(t *testing.T) {
	h := makeHandleJobNextRun(newFakeNextRunWatcher())
	req := httptest.NewRequest(http.MethodGet, "/jobs/nextrun?name=ghost", nil)
	rr := httptest.NewRecorder()
	h(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestHandleJobNextRun_MethodNotAllowed(t *testing.T) {
	h := makeHandleJobNextRun(newFakeNextRunWatcher())
	req := httptest.NewRequest(http.MethodPost, "/jobs/nextrun?name=backup", nil)
	rr := httptest.NewRecorder()
	h(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestHandleJobNextRun_OverdueFlag_SetWhenPast(t *testing.T) {
	fw := newFakeNextRunWatcher()
	fw.entries["cleanup"] = time.Now().Add(-5 * time.Minute)

	h := makeHandleJobNextRun(fw)
	req := httptest.NewRequest(http.MethodGet, "/jobs/nextrun?name=cleanup", nil)
	rr := httptest.NewRecorder()
	h(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&body)
	if body["overdue"] != true {
		t.Errorf("expected overdue=true, got %v", body["overdue"])
	}
}
