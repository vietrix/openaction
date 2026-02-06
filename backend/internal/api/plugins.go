package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) handlePlugins(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,name,description,created_at FROM plugins ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name, desc string
		var created int64
		if err := rows.Scan(&id, &name, &desc, &created); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":          id,
			"name":        name,
			"description": desc,
			"created_at":  created,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreatePlugin(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	id := randomID()
	now := time.Now().Unix()
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO plugins(id,name,description,created_at) VALUES(?,?,?,?)",
		id, payload.Name, payload.Description, now)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "plugins.create", payload.Name, "created", requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handlePluginVersions(w http.ResponseWriter, r *http.Request) {
	pluginID := chiURLParam(r, "id")
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,version,wasm_path,created_at FROM plugin_versions WHERE plugin_id = ? ORDER BY created_at DESC", pluginID)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, version, wasmPath string
		var created int64
		if err := rows.Scan(&id, &version, &wasmPath, &created); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":         id,
			"version":    version,
			"wasm_path":  wasmPath,
			"created_at": created,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreatePluginVersion(w http.ResponseWriter, r *http.Request) {
	pluginID := chiURLParam(r, "id")
	var payload struct {
		Version  string `json:"version"`
		WASMPath string `json:"wasm_path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Version == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.WASMPath == "" {
		payload.WASMPath = "plugins/" + payload.Version + ".wasm"
	}
	id := randomID()
	now := time.Now().Unix()
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO plugin_versions(id,plugin_id,version,wasm_path,created_at) VALUES(?,?,?,?,?)",
		id, pluginID, payload.Version, payload.WASMPath, now)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "plugins.version", pluginID, payload.Version, requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}
