package api

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"openaction/internal/auth"
)

func identityFromContext(r *http.Request) *auth.Identity {
	if r == nil {
		return nil
	}
	if value := r.Context().Value(auth.AuthContextKey{}); value != nil {
		if id, ok := value.(*auth.Identity); ok {
			return id
		}
	}
	return nil
}

func (s *Server) requirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := identityFromContext(r)
			if id == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if id.IsToken {
				next.ServeHTTP(w, r)
				return
			}
			if s.hasPermission(r.Context(), id.UserID, permission) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

func (s *Server) hasPermission(ctx context.Context, userID, permission string) bool {
	var count int
	err := s.DB.QueryRowContext(ctx, `
    SELECT COUNT(1)
    FROM user_roles ur
    JOIN role_permissions rp ON rp.role_id = ur.role_id
    JOIN permissions p ON p.id = rp.permission_id
    WHERE ur.user_id = ? AND (p.name = ? OR p.name = '*')`,
		userID, permission).Scan(&count)
	return err == nil && count > 0
}

func (s *Server) audit(ctx context.Context, actorID, action, resource, payload, ip string) {
	if actorID == "" {
		actorID = "system"
	}
	_, _ = s.DB.ExecContext(ctx, `
    INSERT INTO audit_trail(id,actor_id,action,resource,payload,created_at,ip)
    VALUES(?,?,?,?,?,?,?)`,
		randomID(), actorID, action, resource, payload, time.Now().Unix(), ip)
}

func requestIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func randomID() string {
	return uuid.NewString()
}

func chiURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func identityID(r *http.Request) string {
	if id := identityFromContext(r); id != nil {
		return id.UserID
	}
	return ""
}
