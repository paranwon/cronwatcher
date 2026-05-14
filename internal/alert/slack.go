package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackAlerter sends alerts to a Slack webhook URL.
type SlackAlerter struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlack creates a new SlackAlerter with the given webhook URL.
func NewSlack(webhookURL string) *SlackAlerter {
	return &SlackAlerter{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// MissedJob sends a Slack notification when a cron job is missed.
func (s *SlackAlerter) MissedJob(jobName string, expectedAt time.Time) error {
	msg := fmt.Sprintf(":warning: *Missed cron job*: `%s` was expected at %s",
		jobName, expectedAt.Format(time.RFC3339))
	return s.send(msg)
}

// LongRunningJob sends a Slack notification when a job exceeds its duration threshold.
func (s *SlackAlerter) LongRunningJob(jobName string, duration time.Duration) error {
	msg := fmt.Sprintf(":hourglass_flowing_sand: *Long-running cron job*: `%s` has been running for %s",
		jobName, duration.Round(time.Second))
	return s.send(msg)
}

func (s *SlackAlerter) send(text string) error {
	payload, err := json.Marshal(slackPayload{Text: text})
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("slack: post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
