package watcher

// AlertStateAccessor exposes alert state management on the Watcher.

// HasAlertFired reports whether the given alertType has been fired for job.
func (w *Watcher) HasAlertFired(job, alertType string) bool {
	return w.alertState.HasFired(job, alertType)
}

// SetAlertFired marks alertType as fired for job, suppressing future duplicate alerts.
func (w *Watcher) SetAlertFired(job, alertType string) {
	w.alertState.SetFired(job, alertType)
}

// ClearAlert resets a specific alertType for job, allowing it to fire again.
func (w *Watcher) ClearAlert(job, alertType string) {
	w.alertState.Clear(job, alertType)
}

// ClearAllAlerts resets all alert states for job (e.g. when a job succeeds).
func (w *Watcher) ClearAllAlerts(job string) {
	w.alertState.ClearAll(job)
}
