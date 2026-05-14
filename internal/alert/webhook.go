package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Webhook sends alert notifications to a generic HTTP webhook endpoint.
type Webhook struct {
	url    string
	client *http.Client
}

type webhookPayload struct {
	Event   string `json:"event"`
	Job     string `json:"job"`
	Message string `json:"message"`
	At      string `json:"at"`
}

// NewWebhook creates a new Webhook alerter targeting the given URL.
func NewWebhook(url string) *Webhook {
	return &Webhook{
		url: url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// MissedJob sends a webhook notification for a missed cron job.
func (w *Webhook) MissedJob(jobName string, expectedAt time.Time) error {
	payload := webhookPayload{
		Event:   "missed_job",
		Job:     jobName,
		Message: fmt.Sprintf("Job '%s' did not run as expected.", jobName),
		At:      expectedAt.UTC().Format(time.RFC3339),
	}
	return w.send(payload)
}

// LongRunningJob sends a webhook notification for a job that exceeded its deadline.
func (w *Webhook) LongRunningJob(jobName string, duration time.Duration) error {
	payload := webhookPayload{
		Event:   "long_running_job",
		Job:     jobName,
		Message: fmt.Sprintf("Job '%s' has been running for %s.", jobName, duration.Round(time.Second)),
		At:      time.Now().UTC().Format(time.RFC3339),
	}
	return w.send(payload)
}

func (w *Webhook) send(payload webhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
