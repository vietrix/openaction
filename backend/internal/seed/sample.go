package seed

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"openaction/internal/blob"
	"openaction/internal/db"
)

type sampleProject struct {
	id            string
	name          string
	repoURL       string
	defaultBranch string
}

type samplePipeline struct {
	id          string
	projectID   string
	status      string
	commitHash  string
	branch      string
	triggeredBy string
	startedAt   int64
	finishedAt  int64
}

type sampleStep struct {
	id         string
	pipelineID string
	name       string
	status     string
	startedAt  int64
	finishedAt int64
	logPath    string
}

type sampleRelease struct {
	id         string
	projectID  string
	version    string
	build      string
	patch      string
	createdAt  int64
	updatePath string
	updateMD   string
}

type sampleArtifact struct {
	id        string
	releaseID string
	filename  string
	content   string
	createdAt int64
}

func EnsureSampleData(ctx context.Context, database *db.DB, blobStore *blob.Store, dataDir string) error {
	var count int
	if err := database.QueryRowContext(ctx, "SELECT COUNT(1) FROM projects").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now().Unix()

	projects := []sampleProject{
		{
			id:            uuid.NewString(),
			name:          "forge-api",
			repoURL:       "https://github.com/vietrix/forge-api",
			defaultBranch: "main",
		},
		{
			id:            uuid.NewString(),
			name:          "forge-web",
			repoURL:       "https://github.com/vietrix/forge-web",
			defaultBranch: "main",
		},
	}

	for _, project := range projects {
		if _, err := database.ExecContext(ctx,
			"INSERT INTO projects(id,name,repo_url,default_branch,created_at) VALUES(?,?,?,?,?)",
			project.id, project.name, project.repoURL, project.defaultBranch, now); err != nil {
			return err
		}
	}

	pipelines := []samplePipeline{
		{
			id:          uuid.NewString(),
			projectID:   projects[0].id,
			status:      "running",
			commitHash:  "d4f6a2c",
			branch:      "main",
			triggeredBy: "minh",
			startedAt:   now - 240,
			finishedAt:  0,
		},
		{
			id:          uuid.NewString(),
			projectID:   projects[0].id,
			status:      "success",
			commitHash:  "8ab13ef",
			branch:      "release/0.1",
			triggeredBy: "linh",
			startedAt:   now - 1800,
			finishedAt:  now - 1500,
		},
		{
			id:          uuid.NewString(),
			projectID:   projects[1].id,
			status:      "error",
			commitHash:  "6e20da1",
			branch:      "feature/billing",
			triggeredBy: "bot",
			startedAt:   now - 3600,
			finishedAt:  now - 3300,
		},
	}

	for _, pipeline := range pipelines {
		var finished any = nil
		if pipeline.finishedAt != 0 {
			finished = pipeline.finishedAt
		}
		if _, err := database.ExecContext(ctx, `
      INSERT INTO pipelines(id,project_id,status,commit_hash,branch,triggered_by,started_at,finished_at)
      VALUES(?,?,?,?,?,?,?,?)`,
			pipeline.id, pipeline.projectID, pipeline.status, pipeline.commitHash, pipeline.branch, pipeline.triggeredBy, pipeline.startedAt, finished); err != nil {
			return err
		}
	}

	stepTemplates := []struct {
		name   string
		status string
	}{
		{"checkout", "success"},
		{"install", "success"},
		{"lint", "success"},
		{"test", "running"},
		{"build", "pending"},
		{"deploy", "pending"},
	}

	for _, pipeline := range pipelines {
		var steps []sampleStep
		stepStart := pipeline.startedAt
		for index, tmpl := range stepTemplates {
			status := tmpl.status
			if pipeline.status == "success" {
				status = "success"
			} else if pipeline.status == "error" {
				if tmpl.name == "test" {
					status = "error"
				} else if status == "pending" {
					status = "pending"
				} else {
					status = "success"
				}
			}
			stepID := uuid.NewString()
			startedAt := int64(0)
			finishedAt := int64(0)
			if status != "pending" {
				startedAt = stepStart + int64(index*12)
			}
			if status == "success" {
				finishedAt = startedAt + 10
			}
			if status == "error" {
				finishedAt = startedAt + 6
			}

			logPath := ""
			if status != "pending" {
				content := sampleLogContent(pipeline, tmpl.name, status)
				relPath := blobStore.RelativePath("logs", fmt.Sprintf("%s-%s", pipeline.id, tmpl.name))
				if _, _, err := blobStore.WriteCompressed(relPath, strings.NewReader(content)); err != nil {
					return err
				}
				logPath = relPath
			}

			steps = append(steps, sampleStep{
				id:         stepID,
				pipelineID: pipeline.id,
				name:       tmpl.name,
				status:     status,
				startedAt:  startedAt,
				finishedAt: finishedAt,
				logPath:    logPath,
			})
		}

		for _, step := range steps {
			var started any = nil
			var finished any = nil
			var logPath any = nil
			if step.startedAt != 0 {
				started = step.startedAt
			}
			if step.finishedAt != 0 {
				finished = step.finishedAt
			}
			if step.logPath != "" {
				logPath = step.logPath
			}
			if _, err := database.ExecContext(ctx, `
        INSERT INTO pipeline_steps(id,pipeline_id,name,status,started_at,finished_at,log_path)
        VALUES(?,?,?,?,?,?,?)`,
				step.id, step.pipelineID, step.name, step.status, started, finished, logPath); err != nil {
				return err
			}
		}
	}

	releases := []sampleRelease{
		{
			id:         uuid.NewString(),
			projectID:  projects[0].id,
			version:    "v0.0.3",
			build:      "dev",
			patch:      "a",
			createdAt:  now - 7200,
			updatePath: filepath.ToSlash(filepath.Join("updates", "v0.0.3", "dev", "a.md")),
			updateMD:   sampleUpdateMD("v0.0.3", "dev", "a"),
		},
		{
			id:         uuid.NewString(),
			projectID:  projects[0].id,
			version:    "v0.0.3",
			build:      "prod",
			patch:      "v0.0.3b",
			createdAt:  now - 3600,
			updatePath: filepath.ToSlash(filepath.Join("updates", "v0.0.3", "prod", "v0.0.3b.md")),
			updateMD:   sampleUpdateMD("v0.0.3", "prod", "v0.0.3b"),
		},
		{
			id:         uuid.NewString(),
			projectID:  projects[1].id,
			version:    "bixie",
			build:      "prod",
			patch:      "a",
			createdAt:  now - 1800,
			updatePath: filepath.ToSlash(filepath.Join("updates", "bixie", "prod", "a.md")),
			updateMD:   sampleUpdateMD("bixie", "prod", "a"),
		},
	}

	for _, release := range releases {
		fullPath := filepath.Join(dataDir, release.updatePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err == nil {
			_ = os.WriteFile(fullPath, []byte(release.updateMD), 0o644)
		}
		if _, err := database.ExecContext(ctx, `
      INSERT INTO releases(id,project_id,version,build,patch,created_at,update_path)
      VALUES(?,?,?,?,?,?,?)`,
			release.id, release.projectID, release.version, release.build, release.patch, release.createdAt, release.updatePath); err != nil {
			return err
		}
	}

	artifacts := []sampleArtifact{
		{
			id:        uuid.NewString(),
			releaseID: releases[0].id,
			filename:  "openaction-dev-a-linux-amd64.tar.gz",
			content:   "linux amd64 build content",
			createdAt: now - 7200,
		},
		{
			id:        uuid.NewString(),
			releaseID: releases[0].id,
			filename:  "openaction-dev-a-darwin-arm64.tar.gz",
			content:   "darwin arm64 build content",
			createdAt: now - 7200,
		},
		{
			id:        uuid.NewString(),
			releaseID: releases[1].id,
			filename:  "openaction-prod-v0.0.3b-windows-x64.zip",
			content:   "windows x64 build content",
			createdAt: now - 3600,
		},
		{
			id:        uuid.NewString(),
			releaseID: releases[2].id,
			filename:  "openaction-bixie-prod-a-linux-amd64.tar.gz",
			content:   "bixie prod linux build content",
			createdAt: now - 1800,
		},
	}

	for _, artifact := range artifacts {
		relPath := blobStore.RelativePath("artifacts", fmt.Sprintf("%s-%s", artifact.releaseID, artifact.filename))
		_, size, err := blobStore.WriteCompressed(relPath, strings.NewReader(artifact.content))
		if err != nil {
			return err
		}
		if _, err := database.ExecContext(ctx, `
      INSERT INTO artifacts(id,release_id,filename,size_bytes,blob_path,created_at)
      VALUES(?,?,?,?,?,?)`,
			artifact.id, artifact.releaseID, artifact.filename, size, relPath, artifact.createdAt); err != nil {
			return err
		}
	}

	return nil
}

func sampleLogContent(pipeline samplePipeline, stepName, status string) string {
	lines := []string{
		fmt.Sprintf("→ Pipeline %s (%s) step %s", pipeline.id, pipeline.branch, stepName),
		fmt.Sprintf("[%s] Starting...", stepName),
	}
	switch status {
	case "success":
		lines = append(lines,
			fmt.Sprintf("[%s] Running checks...", stepName),
			fmt.Sprintf("[%s] SUCCESS - Completed", stepName),
		)
	case "error":
		lines = append(lines,
			fmt.Sprintf("[%s] ERROR - Step failed", stepName),
			fmt.Sprintf("[%s] Please inspect logs for details", stepName),
		)
	case "running":
		lines = append(lines,
			fmt.Sprintf("[%s] Streaming output...", stepName),
			fmt.Sprintf("[%s] Still running...", stepName),
		)
	default:
		lines = append(lines, fmt.Sprintf("[%s] Pending", stepName))
	}
	return strings.Join(lines, "\n")
}

func sampleUpdateMD(version, build, patch string) string {
	return fmt.Sprintf(`# Release %s (%s - %s)

## Highlights
- Tăng tốc pipeline với cache thông minh
- Nâng cấp UI theo phong cách Linear
- Cải thiện trải nghiệm log streaming 60fps

## Fixes
- Khắc phục lỗi timeout runner
- Ổn định WASM plugin sandbox
`, version, build, patch)
}
