package alert

import (
	"log/slog"
	"time"

	"github.com/user/cronwatcher/internal/config"
)

// Logger is an AlertSink that writes structured log entries via slog.
type Logger struct {
	log *slog.Logger
}

// NewLogger returns a Logger wrapping the provided slog.Logger.
func NewLogger(l *slog.Logger) *Logger {
	return &Logger{log: l}
}

// MissedJob logs a warning when a job has not been seen within its MaxDelay.
func (l *Logger) MissedJob(job config.Job, lastSeen time.Time) {
	l.log.Warn("missed job",
		slog.String("job", job.Name),
		slog.Time("last_seen", lastSeen),
		slog.Duration("max_delay", job.MaxDelay),
	)
}

// LongRunningJob logs a warning when a job exceeds its MaxDuration.
func (l *Logger) LongRunningJob(job config.Job, duration time.Duration) {
	l.log.Warn("long-running job",
		slog.String("job", job.Name),
		slog.Duration("running_for", duration),
		slog.Duration("max_duration", job.MaxDuration),
	)
}
