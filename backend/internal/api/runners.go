package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) handleRunners(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,name,status,version,last_seen,created_at FROM runners ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name, status, version string
		var lastSeen, created int64
		if err := rows.Scan(&id, &name, &status, &version, &lastSeen, &created); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		tags := s.runnerTags(r.Context(), id)
		items = append(items, map[string]any{
			"id":         id,
			"name":       name,
			"status":     status,
			"version":    version,
			"last_seen":  lastSeen,
			"created_at": created,
			"tags":       tags,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateRunner(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name    string   `json:"name"`
		Status  string   `json:"status"`
		Version string   `json:"version"`
		Tags    []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.Status == "" {
		payload.Status = "online"
	}
	if payload.Version == "" {
		payload.Version = "v1"
	}
	now := time.Now().Unix()
	id := randomID()
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO runners(id,name,status,version,last_seen,created_at) VALUES(?,?,?,?,?,?)",
		id, payload.Name, payload.Status, payload.Version, now, now)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.storeRunnerTags(r.Context(), id, payload.Tags)
	s.audit(r.Context(), identityID(r), "runners.create", payload.Name, "created", requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleUpdateRunner(w http.ResponseWriter, r *http.Request) {
	id := chiURLParam(r, "id")
	var payload struct {
		Status  string   `json:"status"`
		Version string   `json:"version"`
		Tags    []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.Status == "" {
		payload.Status = "online"
	}
	if payload.Version == "" {
		payload.Version = "v1"
	}
	now := time.Now().Unix()
	_, err := s.DB.ExecContext(r.Context(),
		"UPDATE runners SET status = ?, version = ?, last_seen = ? WHERE id = ?",
		payload.Status, payload.Version, now, id)
	if err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	s.storeRunnerTags(r.Context(), id, payload.Tags)
	s.audit(r.Context(), identityID(r), "runners.update", id, "updated", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleDeleteRunner(w http.ResponseWriter, r *http.Request) {
	id := chiURLParam(r, "id")
	_, err := s.DB.ExecContext(r.Context(), "DELETE FROM runners WHERE id = ?", id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "runners.delete", id, "deleted", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleRunnerSummary(w http.ResponseWriter, r *http.Request) {
	var total, online, busy, offline int
	_ = s.DB.QueryRowContext(r.Context(), "SELECT COUNT(1) FROM runners").Scan(&total)
	_ = s.DB.QueryRowContext(r.Context(), "SELECT COUNT(1) FROM runners WHERE status IN ('online','busy')").Scan(&online)
	_ = s.DB.QueryRowContext(r.Context(), "SELECT COUNT(1) FROM runners WHERE status = 'busy'").Scan(&busy)
	_ = s.DB.QueryRowContext(r.Context(), "SELECT COUNT(1) FROM runners WHERE status = 'offline'").Scan(&offline)
	writeJSON(w, http.StatusOK, map[string]any{
		"total":   total,
		"online":  online,
		"busy":    busy,
		"offline": offline,
	})
}

func (s *Server) runnerTags(ctx context.Context, runnerID string) []string {
	rows, err := s.DB.QueryContext(ctx, "SELECT tag FROM runner_tags WHERE runner_id = ?", runnerID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}

func (s *Server) storeRunnerTags(ctx context.Context, runnerID string, tags []string) {
	_, _ = s.DB.ExecContext(ctx, "DELETE FROM runner_tags WHERE runner_id = ?", runnerID)
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		_, _ = s.DB.ExecContext(ctx, "INSERT INTO runner_tags(id,runner_id,tag) VALUES(?,?,?)",
			randomID(), runnerID, tag)
	}
}
