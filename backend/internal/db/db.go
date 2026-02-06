package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", path))
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(2 * time.Minute)

	return &DB{DB: conn}, nil
}

func (d *DB) Close() error {
	if d == nil || d.DB == nil {
		return nil
	}
	return d.DB.Close()
}

func (d *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return d.DB.PingContext(ctx)
}
