package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) handleSSOProviders(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,name,type,config_json,enabled,created_at FROM sso_providers ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name, typ, config string
		var enabled int
		var created int64
		if err := rows.Scan(&id, &name, &typ, &config, &enabled, &created); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":         id,
			"name":       name,
			"type":       typ,
			"config":     config,
			"enabled":    enabled == 1,
			"created_at": created,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateSSOProvider(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Config string `json:"config"`
		Enable bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	id := randomID()
	now := time.Now().Unix()
	enabled := 0
	if payload.Enable {
		enabled = 1
	}
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO sso_providers(id,name,type,config_json,enabled,created_at) VALUES(?,?,?,?,?,?)",
		id, payload.Name, payload.Type, payload.Config, enabled, now)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "sso.create", payload.Name, "created", requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleUpdateSSOProvider(w http.ResponseWriter, r *http.Request) {
	id := chiURLParam(r, "id")
	var payload struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Config string `json:"config"`
		Enable bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	enabled := 0
	if payload.Enable {
		enabled = 1
	}
	_, err := s.DB.ExecContext(r.Context(),
		"UPDATE sso_providers SET name = ?, type = ?, config_json = ?, enabled = ? WHERE id = ?",
		payload.Name, payload.Type, payload.Config, enabled, id)
	if err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "sso.update", id, "updated", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleDeleteSSOProvider(w http.ResponseWriter, r *http.Request) {
	id := chiURLParam(r, "id")
	_, err := s.DB.ExecContext(r.Context(), "DELETE FROM sso_providers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "sso.delete", id, "deleted", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleSSOAuthorize(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "stub",
		"message": "SSO provider integration is not configured yet",
	})
}
