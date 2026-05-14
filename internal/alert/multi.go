package alert

// Multi fans out alert calls to multiple Alerter implementations.
type Multi struct {
	alerters []Alerter
}

// NewMulti creates a Multi alerter that dispatches to all provided alerters.
func NewMulti(alerters ...Alerter) *Multi {
	return &Multi{alerters: alerters}
}

// MissedJob notifies all registered alerters of a missed job.
func (m *Multi) MissedJob(jobName string) {
	for _, a := range m.alerters {
		a.MissedJob(jobName)
	}
}

// LongRunning notifies all registered alerters of a long-running job.
func (m *Multi) LongRunning(jobName string, durationSeconds float64) {
	for _, a := range m.alerters {
		a.LongRunning(jobName, durationSeconds)
	}
}
