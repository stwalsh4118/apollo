package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/rs/zerolog"
)

const (
	createMigrationsTableSQL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
  id TEXT PRIMARY KEY,
  applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`
	insertMigrationSQL = `INSERT INTO schema_migrations(id) VALUES (?);`
	hasMigrationSQL    = `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE id = ?);`
)

func applyMigrations(ctx context.Context, db *sql.DB, migrationFiles fs.FS, logger zerolog.Logger) error {
	if _, err := db.ExecContext(ctx, createMigrationsTableSQL); err != nil {
		return fmt.Errorf("create migration tracking table: %w", err)
	}

	fileNames, err := migrationFileNames(migrationFiles)
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		applied, err := migrationAlreadyApplied(ctx, db, fileName)
		if err != nil {
			return err
		}

		if applied {
			logger.Debug().Str("migration", fileName).Msg("migration already applied")
			continue
		}

		migrationSQL, err := fs.ReadFile(migrationFiles, fileName)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", fileName, err)
		}

		if err := runSingleMigration(ctx, db, fileName, string(migrationSQL)); err != nil {
			return err
		}

		logger.Info().Str("migration", fileName).Msg("migration applied")
	}

	return nil
}

// runSingleMigration executes a migration script and its tracking record inside
// a single transaction. The modernc.org/sqlite driver executes multi-statement
// SQL atomically within one ExecContext call, so all statements in the script
// are covered by the transaction.
func runSingleMigration(ctx context.Context, db *sql.DB, migrationID, script string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration transaction %s: %w", migrationID, err)
	}

	if _, err := tx.ExecContext(ctx, script); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("execute migration %s: %w", migrationID, err)
	}

	if _, err := tx.ExecContext(ctx, insertMigrationSQL, migrationID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert migration %s into tracking table: %w", migrationID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", migrationID, err)
	}

	return nil
}

func migrationAlreadyApplied(ctx context.Context, db *sql.DB, migrationID string) (bool, error) {
	var exists bool
	if err := db.QueryRowContext(ctx, hasMigrationSQL, migrationID).Scan(&exists); err != nil {
		return false, fmt.Errorf("query applied migration %s: %w", migrationID, err)
	}

	return exists, nil
}

func migrationFileNames(migrationFiles fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(migrationFiles, ".")
	if err != nil {
		return nil, fmt.Errorf("list migration files: %w", err)
	}

	fileNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".sql") {
			fileNames = append(fileNames, name)
		}
	}

	sort.Strings(fileNames)

	if len(fileNames) == 0 {
		return nil, fmt.Errorf("no migration files found")
	}

	return fileNames, nil
}
