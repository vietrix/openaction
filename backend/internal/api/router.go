package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"openaction/internal/auth"
	"openaction/internal/blob"
	"openaction/internal/db"
	"openaction/internal/ws"
)

type Server struct {
	DB         *db.DB
	Auth       *auth.Service
	Blob       *blob.Store
	DataDir    string
	SecureOnly bool
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", s.handleHealth)

	r.Route("/public", func(r chi.Router) {
		r.Get("/releases", s.handlePublicReleases)
		r.Get("/releases/{id}", s.handlePublicRelease)
		r.Get("/artifacts/{id}/download", s.handlePublicArtifactDownload)
		r.Get("/latest/{name}", s.handlePublicLatest)
	})

	r.Route("/actions", func(r chi.Router) {
		r.Post("/auth/login", s.handleLogin)
		r.Post("/auth/logout", s.handleLogout)
		r.Get("/auth/tokens", s.handleListTokens)
		r.Post("/auth/tokens", s.handleCreateToken)
		r.Delete("/auth/tokens/{id}", s.handleDeleteToken)

		r.Group(func(r chi.Router) {
			r.Use(s.Auth.Middleware)
			r.Use(s.Auth.CSRFMiddleware)

			r.Get("/projects", s.handleProjects)
			r.Post("/projects", s.handleCreateProject)
			r.Get("/projects/{id}", s.handleProject)
			r.Get("/projects/{id}/pipelines", s.handleProjectPipelines)
			r.Post("/projects/{id}/pipelines", s.handleCreatePipeline)
			r.Get("/pipelines/{id}", s.handlePipeline)
			r.Get("/pipelines/{id}/steps", s.handlePipelineSteps)
			r.Get("/pipelines/{id}/logs", s.handlePipelineLogs)
			r.Get("/pipelines/{id}/logs/stream", wsHandler(s))

			r.Get("/releases", s.handleReleases)
			r.Post("/releases", s.handleCreateRelease)
			r.Get("/releases/{id}", s.handleRelease)
			r.Get("/releases/{id}/artifacts", s.handleArtifacts)
			r.Post("/artifacts", s.handleCreateArtifact)

			r.Get("/settings", s.handleSettings)
			r.Post("/settings", s.handleSettingsUpdate)
			r.Get("/secrets", s.handleSecrets)
			r.Post("/secrets", s.handleSecretsUpdate)
		})
	})

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	user, err := s.Auth.AuthenticateUser(r.Context(), payload.Email, payload.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	sessionID, expires, err := s.Auth.CreateSession(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "cannot create session", http.StatusInternalServerError)
		return
	}
	s.Auth.SetSessionCookie(w, sessionID, expires, s.SecureOnly)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s.Auth.ClearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListTokens(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,name,last_used,created_at FROM api_tokens ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var tokens []map[string]any
	for rows.Next() {
		var id, name string
		var lastUsed, createdAt int64
		if err := rows.Scan(&id, &name, &lastUsed, &createdAt); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		tokens = append(tokens, map[string]any{
			"id":         id,
			"name":       name,
			"last_used":  lastUsed,
			"created_at": createdAt,
		})
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (s *Server) handleCreateToken(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	token, err := s.Auth.CreateToken(r.Context(), payload.Name)
	if err != nil {
		http.Error(w, "cannot create token", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (s *Server) handleDeleteToken(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	_, err := s.DB.ExecContext(r.Context(), "DELETE FROM api_tokens WHERE id = ?", id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), "SELECT id,name,repo_url,default_branch,created_at FROM projects ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name, repo, branch string
		var created int64
		if err := rows.Scan(&id, &name, &repo, &branch, &created); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":             id,
			"name":           name,
			"repo_url":       repo,
			"default_branch": branch,
			"created_at":     created,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name          string `json:"name"`
		RepoURL       string `json:"repo_url"`
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.DefaultBranch == "" {
		payload.DefaultBranch = "main"
	}
	id := uuid.NewString()
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO projects(id,name,repo_url,default_branch,created_at) VALUES(?,?,?,?,?)",
		id, payload.Name, payload.RepoURL, payload.DefaultBranch, time.Now().Unix())
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (s *Server) handleProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var name, repo, branch string
	var created int64
	err := s.DB.QueryRowContext(r.Context(),
		"SELECT id,name,repo_url,default_branch,created_at FROM projects WHERE id = ?", id).
		Scan(&id, &name, &repo, &branch, &created)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":             id,
		"name":           name,
		"repo_url":       repo,
		"default_branch": branch,
		"created_at":     created,
	})
}

func (s *Server) handleProjectPipelines(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	rows, err := s.DB.QueryContext(r.Context(), `
    SELECT id,status,commit_hash,branch,triggered_by,started_at,finished_at
    FROM pipelines WHERE project_id = ? ORDER BY started_at DESC`, projectID)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, status, commit, branch, triggered string
		var started, finished sql.NullInt64
		if err := rows.Scan(&id, &status, &commit, &branch, &triggered, &started, &finished); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":           id,
			"status":       status,
			"commit_hash":  commit,
			"branch":       branch,
			"triggered_by": triggered,
			"started_at":   started.Int64,
			"finished_at":  finished.Int64,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreatePipeline(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	var payload struct {
		CommitHash  string `json:"commit_hash"`
		Branch      string `json:"branch"`
		TriggeredBy string `json:"triggered_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.Branch == "" {
		payload.Branch = "main"
	}
	if payload.TriggeredBy == "" {
		payload.TriggeredBy = "manual"
	}
	id := uuid.NewString()
	_, err := s.DB.ExecContext(r.Context(), `
    INSERT INTO pipelines(id,project_id,status,commit_hash,branch,triggered_by,started_at)
    VALUES(?,?,?,?,?,?,?)`,
		id, projectID, "queued", payload.CommitHash, payload.Branch, payload.TriggeredBy, time.Now().Unix())
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (s *Server) handlePipeline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var projectID, status, commit, branch, triggered string
	var started, finished sql.NullInt64
	err := s.DB.QueryRowContext(r.Context(), `
    SELECT project_id,status,commit_hash,branch,triggered_by,started_at,finished_at
    FROM pipelines WHERE id = ?`, id).
		Scan(&projectID, &status, &commit, &branch, &triggered, &started, &finished)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":           id,
		"project_id":   projectID,
		"status":       status,
		"commit_hash":  commit,
		"branch":       branch,
		"triggered_by": triggered,
		"started_at":   started.Int64,
		"finished_at":  finished.Int64,
	})
}

func (s *Server) handlePipelineSteps(w http.ResponseWriter, r *http.Request) {
	pipelineID := chi.URLParam(r, "id")
	rows, err := s.DB.QueryContext(r.Context(), `
    SELECT id,name,status,started_at,finished_at,log_path
    FROM pipeline_steps WHERE pipeline_id = ? ORDER BY started_at`, pipelineID)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name, status string
		var started, finished sql.NullInt64
		var logPath sql.NullString
		if err := rows.Scan(&id, &name, &status, &started, &finished, &logPath); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":          id,
			"name":        name,
			"status":      status,
			"started_at":  started.Int64,
			"finished_at": finished.Int64,
			"log_path":    logPath.String,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handlePipelineLogs(w http.ResponseWriter, r *http.Request) {
	stepID := r.URL.Query().Get("step_id")
	if stepID == "" {
		http.Error(w, "missing step_id", http.StatusBadRequest)
		return
	}
	var logPath sql.NullString
	err := s.DB.QueryRowContext(r.Context(), "SELECT log_path FROM pipeline_steps WHERE id = ?", stepID).Scan(&logPath)
	if err != nil || !logPath.Valid {
		http.Error(w, "log not found", http.StatusNotFound)
		return
	}

	reader, err := s.Blob.ReadDecompressed(logPath.String)
	if err != nil {
		http.Error(w, "open log failed", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "text/plain")
	_, _ = io.Copy(w, reader)
}

func (s *Server) handleReleases(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.QueryContext(r.Context(), `
    SELECT id,project_id,version,build,patch,created_at,update_path
    FROM releases ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, projectID, version, build, patch, updatePath string
		var created int64
		if err := rows.Scan(&id, &projectID, &version, &build, &patch, &created, &updatePath); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":          id,
			"project_id":  projectID,
			"version":     version,
			"build":       build,
			"patch":       patch,
			"created_at":  created,
			"update_path": updatePath,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateRelease(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ProjectID string `json:"project_id"`
		Version   string `json:"version"`
		Build     string `json:"build"`
		Patch     string `json:"patch"`
		UpdateMD  string `json:"update_md"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	id := uuid.NewString()
	updatePath := filepath.Join("updates", payload.Version, payload.Build, payload.Patch+".md")
	if payload.UpdateMD != "" {
		fullPath := filepath.Join(s.DataDir, updatePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err == nil {
			_ = os.WriteFile(fullPath, []byte(payload.UpdateMD), 0o644)
		}
	}
	_, err := s.DB.ExecContext(r.Context(), `
    INSERT INTO releases(id,project_id,version,build,patch,created_at,update_path)
    VALUES(?,?,?,?,?,?,?)`,
		id, payload.ProjectID, payload.Version, payload.Build, payload.Patch, time.Now().Unix(), updatePath)
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (s *Server) handleRelease(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var projectID, version, build, patch, updatePath string
	var created int64
	err := s.DB.QueryRowContext(r.Context(), `
    SELECT project_id,version,build,patch,created_at,update_path
    FROM releases WHERE id = ?`, id).
		Scan(&projectID, &version, &build, &patch, &created, &updatePath)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	updateContent := ""
	if updatePath != "" {
		fullPath := filepath.Join(s.DataDir, updatePath)
		if data, err := os.ReadFile(fullPath); err == nil {
			updateContent = string(data)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":          id,
		"project_id":  projectID,
		"version":     version,
		"build":       build,
		"patch":       patch,
		"created_at":  created,
		"update_path": updatePath,
		"update_md":   updateContent,
	})
}

func (s *Server) handleArtifacts(w http.ResponseWriter, r *http.Request) {
	releaseID := chi.URLParam(r, "id")
	rows, err := s.DB.QueryContext(r.Context(), `
    SELECT id,filename,size_bytes,blob_path,created_at FROM artifacts WHERE release_id = ?`, releaseID)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name, blobPath string
		var size, created int64
		if err := rows.Scan(&id, &name, &size, &blobPath, &created); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id":         id,
			"filename":   name,
			"size_bytes": size,
			"blob_path":  blobPath,
			"created_at": created,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateArtifact(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ReleaseID string `json:"release_id"`
		Name      string `json:"name"`
		Content   string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.Name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}
	relPath := s.Blob.RelativePath("artifacts", payload.ReleaseID+"-"+payload.Name)
	_, size, err := s.Blob.WriteCompressed(relPath, strings.NewReader(payload.Content))
	if err != nil {
		http.Error(w, "write failed", http.StatusInternalServerError)
		return
	}
	id := uuid.NewString()
	_, err = s.DB.ExecContext(r.Context(), `
    INSERT INTO artifacts(id,release_id,filename,size_bytes,blob_path,created_at)
    VALUES(?,?,?,?,?,?)`,
		id, payload.ReleaseID, payload.Name, size, relPath, time.Now().Unix())
	if err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (s *Server) handlePublicReleases(w http.ResponseWriter, r *http.Request) {
	s.handleReleases(w, r)
}

func (s *Server) handlePublicRelease(w http.ResponseWriter, r *http.Request) {
	s.handleRelease(w, r)
}

func (s *Server) handlePublicArtifactDownload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var name, blobPath string
	err := s.DB.QueryRowContext(r.Context(), "SELECT filename,blob_path FROM artifacts WHERE id = ?", id).Scan(&name, &blobPath)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	reader, err := s.Blob.ReadDecompressed(blobPath)
	if err != nil {
		http.Error(w, "open failed", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = io.Copy(w, reader)
}

func (s *Server) handlePublicLatest(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}
	var id string
	err := s.DB.QueryRowContext(r.Context(), `
    SELECT id FROM artifacts WHERE filename = ? ORDER BY created_at DESC LIMIT 1`, name).Scan(&id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/public/artifacts/"+id+"/download", http.StatusFound)
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"version": "openaction",
		"ui":      "svelte",
	})
}

func (s *Server) handleSettingsUpdate(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleSecrets(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []any{})
}

func (s *Server) handleSecretsUpdate(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func wsHandler(s *Server) http.HandlerFunc {
	handler := &ws.LogHandler{
		Open: func(ctx context.Context, logPath string) (io.ReadCloser, error) {
			return s.Blob.ReadDecompressed(logPath)
		},
	}
	return handler.ServeHTTP
}
