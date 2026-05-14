package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const pagerDutyEventsURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDuty sends alerts to PagerDuty via the Events API v2.
type PagerDuty struct {
	integrationKey string
	endpointURL     string
	client          *http.Client
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
}

// NewPagerDuty creates a new PagerDuty alerter.
func NewPagerDuty(integrationKey string) *PagerDuty {
	return &PagerDuty{
		integrationKey: integrationKey,
		endpointURL:     pagerDutyEventsURL,
		client:          &http.Client{Timeout: 10 * time.Second},
	}
}

// MissedJob sends a PagerDuty trigger event for a missed cron job.
func (p *PagerDuty) MissedJob(jobName string, expectedAt time.Time) error {
	summary := fmt.Sprintf("Cron job '%s' missed its scheduled run at %s",
		jobName, expectedAt.Format(time.RFC3339))
	return p.trigger(summary)
}

// LongRunningJob sends a PagerDuty trigger event for a long-running cron job.
func (p *PagerDuty) LongRunningJob(jobName string, duration time.Duration) error {
	summary := fmt.Sprintf("Cron job '%s' has been running for %s (exceeded threshold)",
		jobName, duration.Round(time.Second))
	return p.trigger(summary)
}

func (p *PagerDuty) trigger(summary string) error {
	body := pdPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:  summary,
			Source:   "cronwatcher",
			Severity: "error",
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", err)
	}

	resp, err := p.client.Post(p.endpointURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: send event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("pagerduty: unexpected status %d: %s", resp.StatusCode, bytes.TrimSpace(respBody))
	}
	return nil
}
