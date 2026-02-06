package api

import (
	"bufio"
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
	SecretKey  []byte
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

			r.With(s.requirePermission("projects.read")).Get("/projects", s.handleProjects)
			r.With(s.requirePermission("projects.write")).Post("/projects", s.handleCreateProject)
			r.With(s.requirePermission("projects.read")).Get("/projects/{id}", s.handleProject)
			r.With(s.requirePermission("pipelines.read")).Get("/projects/{id}/pipelines", s.handleProjectPipelines)
			r.With(s.requirePermission("pipelines.write")).Post("/projects/{id}/pipelines", s.handleCreatePipeline)
			r.With(s.requirePermission("pipelines.read")).Get("/pipelines/{id}", s.handlePipeline)
			r.With(s.requirePermission("pipelines.read")).Get("/pipelines/{id}/steps", s.handlePipelineSteps)
			r.With(s.requirePermission("logs.read")).Get("/pipelines/{id}/logs", s.handlePipelineLogs)
			r.With(s.requirePermission("logs.read")).Get("/pipelines/{id}/logs/stream", wsHandler(s))

			r.With(s.requirePermission("releases.read")).Get("/releases", s.handleReleases)
			r.With(s.requirePermission("releases.write")).Post("/releases", s.handleCreateRelease)
			r.With(s.requirePermission("releases.read")).Get("/releases/{id}", s.handleRelease)
			r.With(s.requirePermission("releases.read")).Get("/releases/{id}/artifacts", s.handleArtifacts)
			r.With(s.requirePermission("artifacts.write")).Post("/artifacts", s.handleCreateArtifact)

			r.With(s.requirePermission("settings.read")).Get("/settings", s.handleSettings)
			r.With(s.requirePermission("settings.write")).Post("/settings", s.handleSettingsUpdate)
			r.With(s.requirePermission("secrets.read")).Get("/secrets", s.handleSecrets)
			r.With(s.requirePermission("secrets.write")).Post("/secrets", s.handleSecretsUpdate)
			r.With(s.requirePermission("secrets.write")).Put("/secrets/{id}", s.handleSecretsUpdate)
			r.With(s.requirePermission("secrets.write")).Delete("/secrets/{id}", s.handleSecretsDelete)

			r.With(s.requirePermission("runners.read")).Get("/runners", s.handleRunners)
			r.With(s.requirePermission("runners.write")).Post("/runners", s.handleCreateRunner)
			r.With(s.requirePermission("runners.write")).Put("/runners/{id}", s.handleUpdateRunner)
			r.With(s.requirePermission("runners.write")).Delete("/runners/{id}", s.handleDeleteRunner)
			r.With(s.requirePermission("runners.read")).Get("/runners/summary", s.handleRunnerSummary)

			r.With(s.requirePermission("env.read")).Get("/environments", s.handleEnvironments)
			r.With(s.requirePermission("env.write")).Post("/environments", s.handleCreateEnvironment)
			r.With(s.requirePermission("env.read")).Get("/environments/{id}/releases", s.handleEnvironmentReleases)
			r.With(s.requirePermission("env.write")).Post("/promotions", s.handlePromotion)
			r.With(s.requirePermission("env.write")).Post("/rollbacks", s.handleRollback)

			r.With(s.requirePermission("plugins.read")).Get("/plugins", s.handlePlugins)
			r.With(s.requirePermission("plugins.write")).Post("/plugins", s.handleCreatePlugin)
			r.With(s.requirePermission("plugins.read")).Get("/plugins/{id}/versions", s.handlePluginVersions)
			r.With(s.requirePermission("plugins.write")).Post("/plugins/{id}/versions", s.handleCreatePluginVersion)

			r.With(s.requirePermission("sso.read")).Get("/sso/providers", s.handleSSOProviders)
			r.With(s.requirePermission("sso.write")).Post("/sso/providers", s.handleCreateSSOProvider)
			r.With(s.requirePermission("sso.write")).Put("/sso/providers/{id}", s.handleUpdateSSOProvider)
			r.With(s.requirePermission("sso.write")).Delete("/sso/providers/{id}", s.handleDeleteSSOProvider)
			r.With(s.requirePermission("sso.read")).Get("/sso/authorize", s.handleSSOAuthorize)

			r.With(s.requirePermission("rbac.read")).Get("/roles", s.handleRoles)
			r.With(s.requirePermission("rbac.write")).Post("/roles", s.handleCreateRole)
			r.With(s.requirePermission("rbac.write")).Put("/roles/{id}", s.handleUpdateRole)
			r.With(s.requirePermission("rbac.write")).Delete("/roles/{id}", s.handleDeleteRole)
			r.With(s.requirePermission("rbac.read")).Get("/permissions", s.handlePermissions)
			r.With(s.requirePermission("rbac.write")).Post("/roles/{id}/permissions", s.handleAssignRolePermissions)
			r.With(s.requirePermission("rbac.write")).Post("/users/{id}/roles", s.handleAssignUserRoles)

			r.With(s.requirePermission("audit.read")).Get("/audit", s.handleAudit)
			r.With(s.requirePermission("metrics.read")).Get("/metrics", s.handleMetrics)
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
	s.audit(r.Context(), identityID(r), "projects.create", payload.Name, "created", requestIP(r))
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
	s.audit(r.Context(), identityID(r), "pipelines.create", projectID, id, requestIP(r))
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
	secrets := s.secretValues(r.Context())
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		for _, secret := range secrets {
			if secret != "" {
				line = strings.ReplaceAll(line, secret, "*****")
			}
		}
		_, _ = io.WriteString(w, line+"\n")
	}
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
	s.audit(r.Context(), identityID(r), "releases.create", payload.ProjectID, id, requestIP(r))
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
	s.audit(r.Context(), identityID(r), "artifacts.create", payload.ReleaseID, payload.Name, requestIP(r))
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
	rows, err := s.DB.QueryContext(r.Context(), "SELECT key,value,updated_at FROM settings ORDER BY key")
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var key, value string
		var updated int64
		if err := rows.Scan(&key, &value, &updated); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"key":        key,
			"value":      value,
			"updated_at": updated,
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleSettingsUpdate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Key == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO settings(key,value,updated_at) VALUES(?,?,?) ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at",
		payload.Key, payload.Value, time.Now().Unix())
	if err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), identityID(r), "settings.update", payload.Key, "updated", requestIP(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func wsHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secrets := s.secretValues(r.Context())
		handler := &ws.LogHandler{
			Open: func(ctx context.Context, logPath string) (io.ReadCloser, error) {
				return s.Blob.ReadDecompressed(logPath)
			},
			Mask: func(line string) string {
				for _, secret := range secrets {
					if secret != "" {
						line = strings.ReplaceAll(line, secret, "*****")
					}
				}
				return line
			},
		}
		handler.ServeHTTP(w, r)
	}
}
