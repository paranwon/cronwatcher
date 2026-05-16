package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/densestvoid/cronwatcher/internal/watcher"
)

type fakeCheckpointWatcher struct {
	data map[string][]watcher.CheckpointEntry
}

func newFakeCheckpointWatcher() *fakeCheckpointWatcher {
	return &fakeCheckpointWatcher{data: map[string][]watcher.CheckpointEntry{"backup": {}}}
}

func (f *fakeCheckpointWatcher) RecordCheckpoint(job, name string, meta map[string]string) error {
	if _, ok := f.data[job]; !ok {
		return fmt.Errorf("unknown job %q", job)
	}
	f.data[job] = append(f.data[job], watcher.CheckpointEntry{Name: name, RecordedAt: time.Now(), Meta: meta})
	return nil
}
func (f *fakeCheckpointWatcher) GetCheckpoints(job string) ([]watcher.CheckpointEntry, bool) {
	e, ok := f.data[job]
	return e, ok
}
func (f *fakeCheckpointWatcher) ClearCheckpoints(job string) error {
	if _, ok := f.data[job]; !ok {
		return fmt.Errorf("unknown job %q", job)
	}
	f.data[job] = nil
	return nil
}

func TestHandleJobCheckpoint_GET_ReturnsEntries(t *testing.T) {
	fw := newFakeCheckpointWatcher()
	_ = fw.RecordCheckpoint("backup", "step-1", nil)
	h := makeHandleJobCheckpoint(fw)
	req := httptest.NewRequest(http.MethodGet, "/job/checkpoint?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var entries []watcher.CheckpointEntry
	json.NewDecoder(rw.Body).Decode(&entries)
	if len(entries) != 1 || entries[0].Name != "step-1" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}

func TestHandleJobCheckpoint_POST_RecordsCheckpoint(t *testing.T) {
	fw := newFakeCheckpointWatcher()
	h := makeHandleJobCheckpoint(fw)
	body, _ := json.Marshal(map[string]interface{}{"checkpoint": "db-dump", "meta": map[string]string{"rows": "100"}})
	req := httptest.NewRequest(http.MethodPost, "/job/checkpoint?name=backup", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rw.Code)
	}
	entries, _ := fw.GetCheckpoints("backup")
	if len(entries) != 1 || entries[0].Name != "db-dump" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}

func TestHandleJobCheckpoint_DELETE_ClearsEntries(t *testing.T) {
	fw := newFakeCheckpointWatcher()
	_ = fw.RecordCheckpoint("backup", "step", nil)
	h := makeHandleJobCheckpoint(fw)
	req := httptest.NewRequest(http.MethodDelete, "/job/checkpoint?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rw.Code)
	}
	entries, _ := fw.GetCheckpoints("backup")
	if len(entries) != 0 {
		t.Errorf("expected empty after delete, got %d", len(entries))
	}
}

func TestHandleJobCheckpoint_MissingName_ReturnsBadRequest(t *testing.T) {
	h := makeHandleJobCheckpoint(newFakeCheckpointWatcher())
	req := httptest.NewRequest(http.MethodGet, "/job/checkpoint", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestHandleJobCheckpoint_MethodNotAllowed(t *testing.T) {
	h := makeHandleJobCheckpoint(newFakeCheckpointWatcher())
	req := httptest.NewRequest(http.MethodPut, "/job/checkpoint?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}
