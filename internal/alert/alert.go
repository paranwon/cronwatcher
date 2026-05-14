package alert

// Alerter defines the interface for sending alerts about cron job issues.
type Alerter interface {
	// MissedJob is called when a cron job has not started within its expected window.
	MissedJob(jobName string)
	// LongRunning is called when a cron job has exceeded its maximum allowed duration.
	LongRunning(jobName string, durationSeconds float64)
}
