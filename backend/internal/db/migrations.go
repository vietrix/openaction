package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const schemaTable = "schema_migrations"

func Migrate(ctx context.Context, db *DB, migrationsDir string) error {
	if db == nil || db.DB == nil {
		return errors.New("db is nil")
	}

	if err := ensureSchemaTable(ctx, db.DB); err != nil {
		return err
	}

	files, err := listMigrationFiles(migrationsDir)
	if err != nil {
		return err
	}

	applied, err := appliedMigrations(ctx, db.DB)
	if err != nil {
		return err
	}

	for _, file := range files {
		if applied[file] {
			continue
		}
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			return err
		}
		if err := applyMigration(ctx, db.DB, file, string(content)); err != nil {
			return fmt.Errorf("apply %s: %w", file, err)
		}
	}

	return nil
}

func ensureSchemaTable(ctx context.Context, db *sql.DB) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		version TEXT PRIMARY KEY,
		applied_at INTEGER NOT NULL
	);`, schemaTable)
	_, err := db.ExecContext(ctx, query)
	return err
}

func listMigrationFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".sql") {
			files = append(files, name)
		}
	}
	sort.Strings(files)
	return files, nil
}

func appliedMigrations(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT version FROM %s", schemaTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

func applyMigration(ctx context.Context, db *sql.DB, version string, sqlText string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, sqlText); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s(version, applied_at) VALUES(?, ?)", schemaTable), version, time.Now().Unix()); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
