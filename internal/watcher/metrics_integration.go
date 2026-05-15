package watcher

import "time"

// GetMetrics returns aggregated metrics for the named job.
// Returns false if the job is unknown or has no recorded runs.
func (w *Watcher) GetMetrics(name string) (JobMetrics, bool) {
	return w.metrics.get(name)
}

// AllMetrics returns a snapshot of metrics for every job that has run.
func (w *Watcher) AllMetrics() []JobMetrics {
	return w.metrics.all()
}

// recordMetrics is called by RecordFinish to update the metrics store.
func (w *Watcher) recordMetrics(name string, d time.Duration, errored bool) {
	w.metrics.record(name, d, errored, time.Now())
}
