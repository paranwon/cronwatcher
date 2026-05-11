package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/cronwatcher/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cronwatcher-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
log_level: debug
jobs:
  - name: backup
    schedule: "0 2 * * *"
    max_duration: 10m
    grace_period: 2m
alerts:
  webhook_url: https://hooks.example.com/alert
`
	path := writeTempConfig(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected log_level=debug, got %q", cfg.LogLevel)
	}
	if len(cfg.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(cfg.Jobs))
	}
	if cfg.Jobs[0].MaxDuration != 10*time.Minute {
		t.Errorf("expected max_duration=10m, got %v", cfg.Jobs[0].MaxDuration)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	yaml := `
jobs:
  - name: cleanup
    schedule: "@daily"
`
	path := writeTempConfig(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log_level=info, got %q", cfg.LogLevel)
	}
	if cfg.Jobs[0].MaxDuration != 5*time.Minute {
		t.Errorf("expected default max_duration=5m, got %v", cfg.Jobs[0].MaxDuration)
	}
	if cfg.Jobs[0].GracePeriod != 1*time.Minute {
		t.Errorf("expected default grace_period=1m, got %v", cfg.Jobs[0].GracePeriod)
	}
}

func TestLoad_MissingJobName(t *testing.T) {
	yaml := `
jobs:
  - schedule: "@hourly"
`
	path := writeTempConfig(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing job name, got nil")
	}
}

func TestLoad_NoJobs(t *testing.T) {
	yaml := `log_level: info
`
	path := writeTempConfig(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for empty jobs list, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
