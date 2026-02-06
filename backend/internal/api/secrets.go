package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"openaction/internal/secret"
)

func (s *Server) handleSecrets(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(),
		"SELECT id,name,scope,created_at,updated_at FROM secrets ORDER BY updated_at DESC")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name, scope string
		var created, updated int64
		if err := rows.Scan(&id, &name, &scope, &created, &updated); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":         id,
			"name":       name,
			"scope":      scope,
			"created_at": created,
			"updated_at": updated,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleSecretsUpdate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		Scope string `json:"scope"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.Name == "" || payload.Value == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	if payload.Scope == "" {
		payload.Scope = "global"
	}
	enc, err := secret.Encrypt(s.SecretKey, payload.Value)
	if err != nil {
		http.Error(w, "encrypt failed", http.StatusInternalServerError)
		return
	}

	now := time.Now().Unix()
	if r.Method == http.MethodPut {
		id := chiURLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		_, err = s.DB.ExecContext(r.Context(),
			"UPDATE secrets SET name = ?, value_enc = ?, scope = ?, updated_at = ? WHERE id = ?",
			payload.Name, enc, payload.Scope, now, id)
		if err != nil {
			http.Error(w, "update failed", http.StatusInternalServerError)
			return
		}
		s.audit(r.Context(), identityID(r), "secrets.update", payload.Name, "updated", requestIP(r))
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	id := randomID()
	_, err = s.DB.ExecContext(r.Context(),
		"INSERT INTO secrets(id,name,value_enc,scope,created_at,updated_at) VALUES(?,?,?,?,?,?)",
		id, payload.Name, enc, payload.Scope, now, now)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "secrets.create", payload.Name, "created", requestIP(r))
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleSecretsDelete(w http.ResponseWriter, r *http.Request) {
	id := chiURLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	_, err := s.DB.ExecContext(r.Context(), "DELETE FROM secrets WHERE id = ?", id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "secrets.delete", id, "deleted", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) secretValues(ctx context.Context) []string {
	rows, err := s.DB.QueryContext(ctx, "SELECT value_enc FROM secrets")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var values []string
	for rows.Next() {
		var enc string
		if err := rows.Scan(&enc); err != nil {
			continue
		}
		plain, err := secret.Decrypt(s.SecretKey, enc)
		if err != nil {
			continue
		}
		if plain != "" {
			values = append(values, plain)
		}
	}
	return values
}
