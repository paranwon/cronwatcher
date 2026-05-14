package scheduler

import (
	"testing"
	"time"

	"cronwatcher/internal/alert"
	"cronwatcher/internal/config"
	"cronwatcher/internal/watcher"
)

func testConfig(jobs []config.Job) *config.Config {
	return &config.Config{
		CheckInterval: config.Duration{Duration: 10 * time.Second},
		Jobs:          jobs,
	}
}

func testWatcher(cfg *config.Config) *watcher.Watcher {
	alerter := alert.NewMulti()
	return watcher.New(cfg, alerter)
}

func TestNew_ReturnsScheduler(t *testing.T) {
	cfg := testConfig(nil)
	w := testWatcher(cfg)
	s := New(cfg, w)
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestRegister_ValidSchedule_NoError(t *testing.T) {
	cfg := testConfig([]config.Job{
		{Name: "heartbeat", Schedule: "@every 1m"},
	})
	w := testWatcher(cfg)
	s := New(cfg, w)

	if err := s.Register(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRegister_InvalidSchedule_ReturnsError(t *testing.T) {
	cfg := testConfig([]config.Job{
		{Name: "bad-job", Schedule: "not-a-cron"},
	})
	w := testWatcher(cfg)
	s := New(cfg, w)

	if err := s.Register(); err == nil {
		t.Fatal("expected error for invalid schedule, got nil")
	}
}

func TestStartStop_DoesNotPanic(t *testing.T) {
	cfg := testConfig([]config.Job{
		{Name: "noop", Schedule: "@every 10m"},
	})
	w := testWatcher(cfg)
	s := New(cfg, w)

	if err := s.Register(); err != nil {
		t.Fatalf("register: %v", err)
	}

	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop() // should not panic or deadlock
}
