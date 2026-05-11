package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Job defines a single cron job to monitor.
type Job struct {
	Name        string        `yaml:"name"`
	Schedule    string        `yaml:"schedule"`
	MaxDuration time.Duration `yaml:"max_duration"`
	GracePeriod time.Duration `yaml:"grace_period"`
}

// AlertConfig holds alerting configuration.
type AlertConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Email      string `yaml:"email"`
}

// Config is the top-level configuration for cronwatcher.
type Config struct {
	LogLevel string      `yaml:"log_level"`
	LogFile  string      `yaml:"log_file"`
	Jobs     []Job       `yaml:"jobs"`
	Alerts   AlertConfig `yaml:"alerts"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate checks required fields and applies sensible defaults.
func (c *Config) validate() error {
	if len(c.Jobs) == 0 {
		return fmt.Errorf("at least one job must be defined")
	}
	for i, job := range c.Jobs {
		if job.Name == "" {
			return fmt.Errorf("job[%d]: name is required", i)
		}
		if job.Schedule == "" {
			return fmt.Errorf("job %q: schedule is required", job.Name)
		}
		if c.Jobs[i].MaxDuration == 0 {
			c.Jobs[i].MaxDuration = 5 * time.Minute
		}
		if c.Jobs[i].GracePeriod == 0 {
			c.Jobs[i].GracePeriod = 1 * time.Minute
		}
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	return nil
}
