package database

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sean/apollo/api/internal/logging"
	"github.com/sean/apollo/api/migrations"
)

func TestOpenRunsMigrationsAndCreatesSchema(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "apollo.db")

	logger, err := logging.New(io.Discard, "info")
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}

	handle, err := Open(ctx, databasePath, logger)
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = handle.Close()
	})

	if _, err := os.Stat(databasePath); err != nil {
		t.Fatalf("expected database file to exist at %s: %v", databasePath, err)
	}

	assertTableExists(t, ctx, handle, "schema_migrations")
	assertTableExists(t, ctx, handle, "topics")
	assertTableExists(t, ctx, handle, "modules")
	assertTableExists(t, ctx, handle, "lessons")
	assertTableExists(t, ctx, handle, "concepts")
	assertTableExists(t, ctx, handle, "concept_references")
	assertTableExists(t, ctx, handle, "topic_prerequisites")
	assertTableExists(t, ctx, handle, "topic_relations")
	assertTableExists(t, ctx, handle, "expansion_queue")
	assertTableExists(t, ctx, handle, "research_jobs")
	assertTableExists(t, ctx, handle, "learning_progress")
	assertTableExists(t, ctx, handle, "concept_retention")
	assertTableExists(t, ctx, handle, "search_index")

	if _, err := handle.DB.ExecContext(ctx, `
		INSERT INTO topics (id, title, status, tags)
		VALUES (?, ?, ?, ?)
	`, "go-basics", "Go Basics", "published", `["go","backend"]`); err != nil {
		t.Fatalf("insert topic with JSON tags: %v", err)
	}

	var firstTag string
	if err := handle.DB.QueryRowContext(ctx, `
		SELECT json_extract(tags, '$[0]')
		FROM topics
		WHERE id = ?
	`, "go-basics").Scan(&firstTag); err != nil {
		t.Fatalf("query json_extract(tags): %v", err)
	}

	if firstTag != "go" {
		t.Fatalf("expected first tag to be %q, got %q", "go", firstTag)
	}
}

func TestMigrationsAreIdempotent(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "apollo.db")

	logger, err := logging.New(io.Discard, "info")
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}

	handle, err := Open(ctx, databasePath, logger)
	if err != nil {
		t.Fatalf("first Open() returned error: %v", err)
	}
	if err := handle.Close(); err != nil {
		t.Fatalf("close first handle: %v", err)
	}

	handle, err = Open(ctx, databasePath, logger)
	if err != nil {
		t.Fatalf("second Open() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = handle.Close()
	})

	fileNames, err := migrationFileNames(migrations.Files)
	if err != nil {
		t.Fatalf("migrationFileNames() returned error: %v", err)
	}

	var appliedCount int
	if err := handle.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&appliedCount); err != nil {
		t.Fatalf("query migration tracking count: %v", err)
	}

	if appliedCount != len(fileNames) {
		t.Fatalf("expected %d applied migrations, got %d", len(fileNames), appliedCount)
	}

	var foreignKeysEnabled int
	if err := handle.DB.QueryRowContext(ctx, `PRAGMA foreign_keys;`).Scan(&foreignKeysEnabled); err != nil {
		t.Fatalf("query PRAGMA foreign_keys: %v", err)
	}

	if foreignKeysEnabled != 1 {
		t.Fatalf("expected foreign keys pragma enabled, got %d", foreignKeysEnabled)
	}
}

func TestOpenFailsOnInvalidPath(t *testing.T) {
	ctx := context.Background()

	logger, err := logging.New(io.Discard, "info")
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}

	_, err = Open(ctx, "/proc/nonexistent/deep/apollo.db", logger)
	if err == nil {
		t.Fatal("expected Open() to return error for invalid path, got nil")
	}
}

func TestCloseOnNilHandle(t *testing.T) {
	var h *Handle
	if err := h.Close(); err != nil {
		t.Fatalf("expected Close() on nil Handle to return nil, got: %v", err)
	}
}

func TestForeignKeysEnforced(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "fk.db")

	logger, err := logging.New(io.Discard, "info")
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}

	handle, err := Open(ctx, databasePath, logger)
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = handle.Close()
	})

	var fkEnabled int
	if err := handle.DB.QueryRowContext(ctx, "PRAGMA foreign_keys;").Scan(&fkEnabled); err != nil {
		t.Fatalf("query PRAGMA foreign_keys: %v", err)
	}

	if fkEnabled != 1 {
		t.Fatalf("expected foreign_keys to be enabled (1), got %d", fkEnabled)
	}

	// Attempt to insert a module referencing a non-existent topic — should fail.
	_, err = handle.DB.ExecContext(ctx, `
		INSERT INTO modules (id, topic_id, title, sort_order)
		VALUES ('mod-1', 'nonexistent-topic', 'Test Module', 1)
	`)
	if err == nil {
		t.Fatal("expected foreign key violation error, got nil")
	}
}

func TestSearchIndexIsFTS5(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "fts.db")

	logger, err := logging.New(io.Discard, "info")
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}

	handle, err := Open(ctx, databasePath, logger)
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = handle.Close()
	})

	var sqlDef string
	err = handle.DB.QueryRowContext(ctx, `
		SELECT sql FROM sqlite_master
		WHERE type = 'table' AND name = 'search_index'
	`).Scan(&sqlDef)
	if err != nil {
		t.Fatalf("query sqlite_master for search_index: %v", err)
	}

	if !strings.Contains(strings.ToLower(sqlDef), "fts5") {
		t.Fatalf("expected search_index to be an FTS5 virtual table, got SQL: %s", sqlDef)
	}
}

func assertTableExists(t *testing.T, ctx context.Context, handle *Handle, tableName string) {
	t.Helper()

	var exists bool
	if err := handle.DB.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM sqlite_master
			WHERE type IN ('table', 'view')
			  AND name = ?
		)
	`, tableName).Scan(&exists); err != nil {
		t.Fatalf("query sqlite_master for table %s: %v", tableName, err)
	}

	if !exists {
		t.Fatalf("expected table %s to exist", tableName)
	}
}
