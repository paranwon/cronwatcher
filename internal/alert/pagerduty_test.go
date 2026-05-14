package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestPagerDutyServer(t *testing.T, statusCode int, capture *pdPayload) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if capture != nil {
			if err := json.NewDecoder(r.Body).Decode(capture); err != nil {
				t.Errorf("failed to decode request body: %v", err)
			}
		}
		w.WriteHeader(statusCode)
	}))
}

func TestPagerDuty_MissedJob_SendsCorrectPayload(t *testing.T) {
	var captured pdPayload
	server := newTestPagerDutyServer(t, http.StatusAccepted, &captured)
	defer server.Close()

	pd := NewPagerDuty("test-integration-key")
	pd.endpointURL = server.URL

	expectedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if err := pd.MissedJob("backup", expectedAt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.RoutingKey != "test-integration-key" {
		t.Errorf("expected routing key %q, got %q", "test-integration-key", captured.RoutingKey)
	}
	if captured.EventAction != "trigger" {
		t.Errorf("expected event_action 'trigger', got %q", captured.EventAction)
	}
	if !strings.Contains(captured.Payload.Summary, "backup") {
		t.Errorf("expected summary to contain job name 'backup', got %q", captured.Payload.Summary)
	}
	if captured.Payload.Severity != "error" {
		t.Errorf("expected severity 'error', got %q", captured.Payload.Severity)
	}
	if captured.Payload.Source != "cronwatcher" {
		t.Errorf("expected source 'cronwatcher', got %q", captured.Payload.Source)
	}
}

func TestPagerDuty_LongRunningJob_SendsCorrectPayload(t *testing.T) {
	var captured pdPayload
	server := newTestPagerDutyServer(t, http.StatusAccepted, &captured)
	defer server.Close()

	pd := NewPagerDuty("key-123")
	pd.endpointURL = server.URL

	if err := pd.LongRunningJob("report-gen", 45*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(captured.Payload.Summary, "report-gen") {
		t.Errorf("expected summary to contain 'report-gen', got %q", captured.Payload.Summary)
	}
	if !strings.Contains(captured.Payload.Summary, "45m") {
		t.Errorf("expected summary to contain duration '45m', got %q", captured.Payload.Summary)
	}
}

func TestPagerDuty_NonOKStatus_ReturnsError(t *testing.T) {
	server := newTestPagerDutyServer(t, http.StatusUnauthorized, nil)
	defer server.Close()

	pd := NewPagerDuty("bad-key")
	pd.endpointURL = server.URL

	err := pd.MissedJob("cleanup", time.Now())
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected error to mention status 401, got: %v", err)
	}
}
