package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestSlackServer(t *testing.T, statusCode int, capturedBody *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf strings.Builder
		_, _ = buf.ReadFrom(r.Body)
		*capturedBody = buf.String()
		w.WriteHeader(statusCode)
	}))
}

func TestSlack_MissedJob_SendsCorrectPayload(t *testing.T) {
	var body string
	server := newTestSlackServer(t, http.StatusOK, &body)
	defer server.Close()

	s := NewSlack(server.URL)
	expectedAt := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	if err := s.MissedJob("backup", expectedAt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload slackPayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("failed to parse payload: %v", err)
	}
	if !strings.Contains(payload.Text, "backup") {
		t.Errorf("expected payload to mention job name, got: %s", payload.Text)
	}
	if !strings.Contains(payload.Text, "Missed") {
		t.Errorf("expected payload to indicate missed job, got: %s", payload.Text)
	}
}

func TestSlack_LongRunningJob_SendsCorrectPayload(t *testing.T) {
	var body string
	server := newTestSlackServer(t, http.StatusOK, &body)
	defer server.Close()

	s := NewSlack(server.URL)

	if err := s.LongRunningJob("report-gen", 5*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload slackPayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("failed to parse payload: %v", err)
	}
	if !strings.Contains(payload.Text, "report-gen") {
		t.Errorf("expected payload to mention job name, got: %s", payload.Text)
	}
}

func TestSlack_NonOKStatus_ReturnsError(t *testing.T) {
	var body string
	server := newTestSlackServer(t, http.StatusInternalServerError, &body)
	defer server.Close()

	s := NewSlack(server.URL)
	err := s.MissedJob("cleanup", time.Now())
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code, got: %v", err)
	}
}
