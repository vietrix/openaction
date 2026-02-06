package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) handleRoles(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,name FROM roles ORDER BY name")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{"id": id, "name": name})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	id := randomID()
	_, err := s.DB.ExecContext(r.Context(), "INSERT INTO roles(id,name) VALUES(?,?)", id, payload.Name)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "roles.create", payload.Name, "created", requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleUpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chiURLParam(r, "id")
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	_, err := s.DB.ExecContext(r.Context(), "UPDATE roles SET name = ? WHERE id = ?", payload.Name, id)
	if err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "roles.update", payload.Name, "updated", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleDeleteRole(w http.ResponseWriter, r *http.Request) {
	id := chiURLParam(r, "id")
	_, err := s.DB.ExecContext(r.Context(), "DELETE FROM roles WHERE id = ?", id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "roles.delete", id, "deleted", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handlePermissions(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,name FROM permissions ORDER BY name")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{"id": id, "name": name})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleAssignRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleID := chiURLParam(r, "id")
	var payload struct {
		Permissions []string `json:"permissions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if len(payload.Permissions) == 0 {
		http.Error(w, "missing permissions", http.StatusBadRequest)
		return
	}
	now := time.Now().Unix()
	for _, perm := range payload.Permissions {
		var permID string
		err := s.DB.QueryRowContext(r.Context(), "SELECT id FROM permissions WHERE name = ?", perm).Scan(&permID)
		if err != nil {
			permID = randomID()
			_, _ = s.DB.ExecContext(r.Context(), "INSERT INTO permissions(id,name) VALUES(?,?)", permID, perm)
		}
		_, _ = s.DB.ExecContext(r.Context(),
			"INSERT INTO role_permissions(id,role_id,permission_id) VALUES(?,?,?)",
			randomID(), roleID, permID)
		_ = now
	}
	s.audit(r.Context(), identityID(r), "roles.permissions", roleID, "updated", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleAssignUserRoles(w http.ResponseWriter, r *http.Request) {
	userID := chiURLParam(r, "id")
	var payload struct {
		Roles []string `json:"roles"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if len(payload.Roles) == 0 {
		http.Error(w, "missing roles", http.StatusBadRequest)
		return
	}
	for _, roleID := range payload.Roles {
		_, _ = s.DB.ExecContext(r.Context(),
			"INSERT INTO user_roles(id,user_id,role_id) VALUES(?,?,?)",
			randomID(), userID, roleID)
	}
	s.audit(r.Context(), identityID(r), "roles.assign", userID, "assigned", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
