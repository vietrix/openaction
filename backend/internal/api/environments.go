package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) handleEnvironments(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,name,created_at FROM environments ORDER BY created_at")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name string
		var created int64
		if err := rows.Scan(&id, &name, &created); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{"id": id, "name": name, "created_at": created})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateEnvironment(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	id := randomID()
	now := time.Now().Unix()
	_, err := s.DB.ExecContext(r.Context(), "INSERT INTO environments(id,name,created_at) VALUES(?,?,?)", id, payload.Name, now)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "env.create", payload.Name, "created", requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleEnvironmentReleases(w http.ResponseWriter, r *http.Request) {
	envID := chiURLParam(r, "id")
	rows, err := s.DB.QueryContext(r.Context(), `
    SELECT environment_releases.id, releases.id, releases.version, releases.build, releases.patch, environment_releases.promoted_at
    FROM environment_releases
    JOIN releases ON releases.id = environment_releases.release_id
    WHERE environment_releases.environment_id = ?
    ORDER BY environment_releases.promoted_at DESC`, envID)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var envRelID, releaseID, version, build, patch string
		var promoted int64
		if err := rows.Scan(&envRelID, &releaseID, &version, &build, &patch, &promoted); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":          envRelID,
			"release_id":  releaseID,
			"version":     version,
			"build":       build,
			"patch":       patch,
			"promoted_at": promoted,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handlePromotion(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		EnvironmentID string `json:"environment_id"`
		ReleaseID     string `json:"release_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.EnvironmentID == "" || payload.ReleaseID == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	now := time.Now().Unix()
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO promotions(id,environment_id,release_id,actor_id,created_at) VALUES(?,?,?,?,?)",
		randomID(), payload.EnvironmentID, payload.ReleaseID, identityID(r), now)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	_, _ = s.DB.ExecContext(r.Context(),
		"INSERT INTO environment_releases(id,environment_id,release_id,promoted_at) VALUES(?,?,?,?)",
		randomID(), payload.EnvironmentID, payload.ReleaseID, now)
	s.audit(r.Context(), identityID(r), "env.promote", payload.EnvironmentID, payload.ReleaseID, requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (s *Server) handleRollback(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		EnvironmentID string `json:"environment_id"`
		ReleaseID     string `json:"release_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.EnvironmentID == "" || payload.ReleaseID == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	now := time.Now().Unix()
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO rollbacks(id,environment_id,release_id,actor_id,created_at) VALUES(?,?,?,?,?)",
		randomID(), payload.EnvironmentID, payload.ReleaseID, identityID(r), now)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	_, _ = s.DB.ExecContext(r.Context(),
		"INSERT INTO environment_releases(id,environment_id,release_id,promoted_at) VALUES(?,?,?,?)",
		randomID(), payload.EnvironmentID, payload.ReleaseID, now)
	s.audit(r.Context(), identityID(r), "env.rollback", payload.EnvironmentID, payload.ReleaseID, requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}
