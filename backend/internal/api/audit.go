package api

import (
	"net/http"
)

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), `
    SELECT id,actor_id,action,resource,payload,created_at,ip
    FROM audit_trail ORDER BY created_at DESC LIMIT 200`)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, actor, action, resource, payload, ip string
		var created int64
		if err := rows.Scan(&id, &actor, &action, &resource, &payload, &created, &ip); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":         id,
			"actor_id":   actor,
			"action":     action,
			"resource":   resource,
			"payload":    payload,
			"created_at": created,
			"ip":         ip,
		})
	}
	writeJSON(w, http.StatusOK, items)
}
