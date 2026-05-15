package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type metricsProvider interface {
	GetMetrics(name string) (interface{ jobMetricsData() jobMetricsResponse }, bool)
}

type jobMetricsResponse struct {
	JobName     string        `json:"job_name"`
	RunCount    int           `json:"run_count"`
	ErrorCount  int           `json:"error_count"`
	AvgDuration string        `json:"avg_duration"`
	MinDuration string        `json:"min_duration"`
	MaxDuration string        `json:"max_duration"`
	LastRun     *time.Time    `json:"last_run,omitempty"`
}

func makeHandleJobMetrics(w statusWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		m, ok := w.GetMetrics(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}

		resp := jobMetricsResponse{
			JobName:     m.JobName,
			RunCount:    m.RunCount,
			ErrorCount:  m.ErrorCount,
			AvgDuration: m.AvgDuration().String(),
			MinDuration: m.MinDuration.String(),
			MaxDuration: m.MaxDuration.String(),
		}
		if !m.LastRun.IsZero() {
			t := m.LastRun
			resp.LastRun = &t
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(resp)
	}
}
