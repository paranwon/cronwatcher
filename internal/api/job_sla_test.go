package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeSLAWatcher struct {
	entries map[string]slaEntry
}

func (f *fakeSLAWatcher) GetSLA(name string) (slaEntry, bool) {
	e, ok := f.entries[name]
	return e, ok
}

func newFakeSLAWatcher() *fakeSLAWatcher {
	return &fakeSLAWatcher{entries: make(map[string]slaEntry)}
}

func TestHandleJobSLA_ReturnsSLAEntry(t *testing.T) {
	fw := newFakeSLAWatcher()
	deadline := time.Now().Add(time.Hour).Truncate(time.Second)
	fw.entries["backup"] = slaEntry{
		JobName:  "backup",
		Deadline: deadline,
		MetSLA:   true,
	}

	h := makeHandleJobSLA(fw)
	req := httptest.NewRequest(http.MethodGet, "/jobs/sla?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var got slaEntry
	if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !got.MetSLA {
		t.Error("expected MetSLA true")
	}
}

func TestHandleJobSLA_MissingName_ReturnsBadRequest(t *testing.T) {
	h := makeHandleJobSLA(newFakeSLAWatcher())
	req := httptest.NewRequest(http.MethodGet, "/jobs/sla", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestHandleJobSLA_NotFound_Returns404(t *testing.T) {
	h := makeHandleJobSLA(newFakeSLAWatcher())
	req := httptest.NewRequest(http.MethodGet, "/jobs/sla?name=ghost", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rw.Code)
	}
}

func TestHandleJobSLA_MethodNotAllowed(t *testing.T) {
	h := makeHandleJobSLA(newFakeSLAWatcher())
	req := httptest.NewRequest(http.MethodPost, "/jobs/sla?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}
