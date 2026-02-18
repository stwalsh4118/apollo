package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/migrations"

	_ "modernc.org/sqlite"
)

const (
	sqliteDriverName     = "sqlite"
	journalModePragmaSQL = "PRAGMA journal_mode = WAL;"
	maxOpenConnections   = 1
	maxIdleConnections   = 1
	databaseDirPerms     = 0750
	// DSN pragma parameters ensure foreign keys and busy timeout are applied to
	// every connection created by database/sql, surviving pool recycling.
	dsnPragmas = "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
)

// Handle wraps a database connection for Apollo services.
type Handle struct {
	DB *sql.DB
}

// Open creates the SQLite database connection, runs migrations, and performs a health check.
func Open(ctx context.Context, databasePath string, logger zerolog.Logger) (*Handle, error) {
	if err := os.MkdirAll(filepath.Dir(databasePath), databaseDirPerms); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := sql.Open(sqliteDriverName, databasePath+dsnPragmas)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)

	// journal_mode=WAL is persistent (written to the DB file) so it only
	// needs to be set once rather than on every new connection.
	if err := applyPersistentPragmas(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := applyMigrations(ctx, db, migrations.Files, logger); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := HealthCheck(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Handle{DB: db}, nil
}

// Close releases database resources.
func (h *Handle) Close() error {
	if h == nil || h.DB == nil {
		return nil
	}

	return h.DB.Close()
}

// HealthCheck validates the database connection.
func HealthCheck(ctx context.Context, db *sql.DB) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping sqlite database: %w", err)
	}

	return nil
}

func applyPersistentPragmas(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, journalModePragmaSQL); err != nil {
		return fmt.Errorf("apply pragma %q: %w", journalModePragmaSQL, err)
	}

	return nil
}
