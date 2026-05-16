package watcher

// RecordCheckpoint records a named checkpoint for the given job.
func (w *Watcher) RecordCheckpoint(job, name string, meta map[string]string) error {
	return w.checkpoints.RecordCheckpoint(job, name, meta)
}

// GetCheckpoints returns all checkpoints recorded for a job.
func (w *Watcher) GetCheckpoints(job string) ([]CheckpointEntry, bool) {
	return w.checkpoints.GetCheckpoints(job)
}

// ClearCheckpoints removes all checkpoints for a job.
func (w *Watcher) ClearCheckpoints(job string) error {
	return w.checkpoints.ClearCheckpoints(job)
}
