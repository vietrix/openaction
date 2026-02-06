package api

import (
	"net/http"
	"time"
)

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	lookback := time.Now().Add(-30 * 24 * time.Hour).Unix()
	var total int
	var failed int
	var running int
	_ = s.DB.QueryRowContext(r.Context(), "SELECT COUNT(1) FROM pipelines WHERE started_at >= ?", lookback).Scan(&total)
	_ = s.DB.QueryRowContext(r.Context(), "SELECT COUNT(1) FROM pipelines WHERE status = 'error' AND started_at >= ?", lookback).Scan(&failed)
	_ = s.DB.QueryRowContext(r.Context(), "SELECT COUNT(1) FROM pipelines WHERE status = 'running' AND started_at >= ?", lookback).Scan(&running)

	var avgDuration float64
	_ = s.DB.QueryRowContext(r.Context(), `
    SELECT AVG(finished_at - started_at)
    FROM pipelines
    WHERE finished_at IS NOT NULL AND started_at >= ?`, lookback).Scan(&avgDuration)

	failureRate := 0.0
	if total > 0 {
		failureRate = float64(failed) / float64(total)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total":        total,
		"failed":       failed,
		"running":      running,
		"avg_duration": avgDuration,
		"failure_rate": failureRate,
		"window_days":  30,
		"generated_at": time.Now().Unix(),
	})
}
