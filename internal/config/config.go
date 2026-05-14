package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Job describes a single cron job to monitor.
type Job struct {
	Name        string        `yaml:"name"`
	Schedule    string        `yaml:"schedule"`
	MaxDelay    time.Duration `yaml:"max_delay"`
	MaxDuration time.Duration `yaml:"max_duration"`
}

// Config is the top-level configuration structure.
type Config struct {
	CheckInterval time.Duration `yaml:"check_interval"`
	LogLevel      string        `yaml:"log_level"`
	Jobs          []Job         `yaml:"jobs"`
}

const (
	defaultCheckInterval = 1 * time.Minute
	defaultLogLevel      = "info"
	defaultMaxDelay      = 10 * time.Minute
	defaultMaxDuration   = 30 * time.Minute
)

// Load reads and validates a YAML config file from path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	applyDefaults(&cfg)

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = defaultCheckInterval
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = defaultLogLevel
	}
	for i := range cfg.Jobs {
		if cfg.Jobs[i].MaxDelay == 0 {
			cfg.Jobs[i].MaxDelay = defaultMaxDelay
		}
		if cfg.Jobs[i].MaxDuration == 0 {
			cfg.Jobs[i].MaxDuration = defaultMaxDuration
		}
	}
}

func validate(cfg *Config) error {
	if len(cfg.Jobs) == 0 {
		return errors.New("config must define at least one job")
	}
	for _, j := range cfg.Jobs {
		if j.Name == "" {
			return errors.New("each job must have a name")
		}
	}
	return nil
}
