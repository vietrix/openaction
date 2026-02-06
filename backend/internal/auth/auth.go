package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Service struct {
	DB         Store
	SessionTTL time.Duration
	TokenTTL   time.Duration
	CSRFFlag   bool
}

type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
}

type AuthContextKey struct{}

type Identity struct {
	UserID  string
	Email   string
	IsToken bool
}

func (s *Service) EnsureAdmin(ctx context.Context, email, password string) error {
	var count int
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM users").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.DB.ExecContext(ctx, "INSERT INTO users(id,email,name,password_hash,created_at) VALUES(?,?,?,?,?)",
		uuid.NewString(), email, "Admin", string(hash), time.Now().Unix())
	return err
}

func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	var user User
	err := s.DB.QueryRowContext(ctx, "SELECT id,email,name,password_hash FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}

func (s *Service) CreateSession(ctx context.Context, userID string) (string, time.Time, error) {
	sessionID := "sess_" + uuid.NewString()
	expires := time.Now().Add(s.SessionTTL)
	_, err := s.DB.ExecContext(ctx, "INSERT INTO sessions(session_id,user_id,expires_at,created_at) VALUES(?,?,?,?)",
		sessionID, userID, expires.Unix(), time.Now().Unix())
	return sessionID, expires, err
}

func (s *Service) CreateToken(ctx context.Context, name string) (string, error) {
	raw := "oa_" + randomString(32)
	hash := sha256.Sum256([]byte(raw))
	expiresAt := time.Now().Add(s.TokenTTL).Unix()
	_, err := s.DB.ExecContext(ctx,
		"INSERT INTO api_tokens(id,name,token_hash,last_used,created_at,expires_at) VALUES(?,?,?,?,?,?)",
		uuid.NewString(), name, hex.EncodeToString(hash[:]), time.Now().Unix(), time.Now().Unix(), expiresAt)
	if err != nil {
		return "", err
	}
	return raw, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (*Identity, error) {
	hash := sha256.Sum256([]byte(token))
	var id string
	var name string
	var expiresAt sql.NullInt64
	err := s.DB.QueryRowContext(ctx, "SELECT id,name,expires_at FROM api_tokens WHERE token_hash = ?", hex.EncodeToString(hash[:])).
		Scan(&id, &name, &expiresAt)
	if err != nil {
		return nil, err
	}
	if expiresAt.Valid && time.Now().Unix() > expiresAt.Int64 {
		_, _ = s.DB.ExecContext(ctx, "DELETE FROM api_tokens WHERE id = ?", id)
		return nil, errors.New("token expired")
	}
	_, _ = s.DB.ExecContext(ctx, "UPDATE api_tokens SET last_used = ? WHERE id = ?", time.Now().Unix(), id)
	return &Identity{UserID: id, Email: name, IsToken: true}, nil
}

func (s *Service) ValidateSession(ctx context.Context, sessionID string) (*Identity, error) {
	var userID, email string
	var expires int64
	err := s.DB.QueryRowContext(ctx, `
    SELECT users.id, users.email, sessions.expires_at
    FROM sessions
    JOIN users ON users.id = sessions.user_id
    WHERE sessions.session_id = ?`, sessionID).Scan(&userID, &email, &expires)
	if err != nil {
		return nil, err
	}
	if time.Now().Unix() > expires {
		_, _ = s.DB.ExecContext(ctx, "DELETE FROM sessions WHERE session_id = ?", sessionID)
		return nil, errors.New("session expired")
	}
	return &Identity{UserID: userID, Email: email, IsToken: false}, nil
}

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := identityFromRequest(r, s); id != nil {
			ctx := context.WithValue(r.Context(), AuthContextKey{}, id)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}

func (s *Service) CSRFMiddleware(next http.Handler) http.Handler {
	if !s.CSRFFlag {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := identityFromRequest(r, s)
		if id == nil || id.IsToken {
			next.ServeHTTP(w, r)
			return
		}
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		csrfCookie, err := r.Cookie("oa_csrf")
		if err != nil {
			http.Error(w, "missing csrf", http.StatusForbidden)
			return
		}
		csrfHeader := r.Header.Get("X-CSRF-Token")
		if csrfHeader == "" || csrfHeader != csrfCookie.Value {
			http.Error(w, "invalid csrf", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Service) SetSessionCookie(w http.ResponseWriter, sessionID string, expires time.Time, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "oa_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
	if s.CSRFFlag {
		csrfValue := randomString(24)
		http.SetCookie(w, &http.Cookie{
			Name:     "oa_csrf",
			Value:    csrfValue,
			Path:     "/",
			Expires:  expires,
			HttpOnly: false,
			Secure:   secure,
			SameSite: http.SameSiteStrictMode,
		})
	}
}

func (s *Service) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "oa_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "oa_csrf",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (s *Service) CleanupExpired(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().Unix()
			_, _ = s.DB.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at < ?", now)
			_, _ = s.DB.ExecContext(ctx, "DELETE FROM api_tokens WHERE expires_at < ?", now)
		}
	}
}

func identityFromRequest(r *http.Request, s *Service) *Identity {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token != "" {
			if id, err := s.ValidateToken(r.Context(), token); err == nil {
				return id
			}
		}
	}

	if cookie, err := r.Cookie("oa_session"); err == nil {
		if id, err := s.ValidateSession(r.Context(), cookie.Value); err == nil {
			return id
		}
	}
	return nil
}

func randomString(length int) string {
	buf := make([]byte, length)
	_, _ = rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)[:length]
}
