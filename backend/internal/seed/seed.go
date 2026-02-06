package seed

import (
	"context"
	"time"

	"github.com/google/uuid"

	"openaction/internal/db"
)

func EnsureDefaults(ctx context.Context, database *db.DB) error {
	if err := ensureEnvironments(ctx, database); err != nil {
		return err
	}
	if err := ensureRBAC(ctx, database); err != nil {
		return err
	}
	return nil
}

func ensureEnvironments(ctx context.Context, database *db.DB) error {
	var count int
	if err := database.QueryRowContext(ctx, "SELECT COUNT(1) FROM environments").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	now := time.Now().Unix()
	for _, name := range []string{"dev", "staging", "prod"} {
		if _, err := database.ExecContext(ctx,
			"INSERT INTO environments(id,name,created_at) VALUES(?,?,?)",
			uuid.NewString(), name, now); err != nil {
			return err
		}
	}
	return nil
}

func ensureRBAC(ctx context.Context, database *db.DB) error {
	var count int
	if err := database.QueryRowContext(ctx, "SELECT COUNT(1) FROM roles").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	adminRoleID := uuid.NewString()
	if _, err := database.ExecContext(ctx, "INSERT INTO roles(id,name) VALUES(?,?)", adminRoleID, "admin"); err != nil {
		return err
	}

	permissions := []string{
		"*",
		"projects.read",
		"projects.write",
		"pipelines.read",
		"pipelines.write",
		"logs.read",
		"releases.read",
		"releases.write",
		"artifacts.write",
		"settings.read",
		"settings.write",
		"secrets.read",
		"secrets.write",
		"runners.read",
		"runners.write",
		"env.read",
		"env.write",
		"plugins.read",
		"plugins.write",
		"sso.read",
		"sso.write",
		"rbac.read",
		"rbac.write",
		"audit.read",
		"metrics.read",
	}
	var wildcardID string
	for _, name := range permissions {
		permID := uuid.NewString()
		if _, err := database.ExecContext(ctx, "INSERT INTO permissions(id,name) VALUES(?,?)", permID, name); err == nil {
			if name == "*" {
				wildcardID = permID
			}
		} else if name == "*" {
			_ = database.QueryRowContext(ctx, "SELECT id FROM permissions WHERE name = '*'").Scan(&wildcardID)
		}
	}

	if wildcardID == "" {
		return nil
	}

	if _, err := database.ExecContext(ctx,
		"INSERT INTO role_permissions(id,role_id,permission_id) VALUES(?,?,?)",
		uuid.NewString(), adminRoleID, wildcardID); err != nil {
		return err
	}

	var adminUserID string
	if err := database.QueryRowContext(ctx, "SELECT id FROM users ORDER BY created_at LIMIT 1").Scan(&adminUserID); err != nil {
		return err
	}
	_, err := database.ExecContext(ctx,
		"INSERT INTO user_roles(id,user_id,role_id) VALUES(?,?,?)",
		uuid.NewString(), adminUserID, adminRoleID)
	return err
}
