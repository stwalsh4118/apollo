package research_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/internal/database"
	"github.com/sean/apollo/api/internal/research"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	logger := zerolog.New(os.Stderr).Level(zerolog.Disabled)

	handle, err := database.Open(context.Background(), dbPath, logger)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() { _ = handle.Close() })

	return handle.DB
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec %q: %v", query, err)
	}
}

func TestPoolSummaryBuildEmpty(t *testing.T) {
	db := setupTestDB(t)
	builder := research.NewPoolSummaryBuilder(db)

	data, err := builder.Build(context.Background())
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	var summary research.PoolSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(summary.ExistingTopics) != 0 {
		t.Fatalf("expected 0 topics, got %d", len(summary.ExistingTopics))
	}

	if len(summary.ExistingConcepts) != 0 {
		t.Fatalf("expected 0 concepts, got %d", len(summary.ExistingConcepts))
	}
}

func TestPoolSummaryBuildPopulated(t *testing.T) {
	db := setupTestDB(t)

	mustExec(t, db, `INSERT INTO topics (id, title, status, tags) VALUES (?, ?, ?, ?)`,
		"go-basics", "Go Basics", "published", `["go"]`)
	mustExec(t, db, `INSERT INTO modules (id, topic_id, title, sort_order) VALUES (?, ?, ?, ?)`,
		"mod-1", "go-basics", "Module 1", 1)
	mustExec(t, db, `INSERT INTO modules (id, topic_id, title, sort_order) VALUES (?, ?, ?, ?)`,
		"mod-2", "go-basics", "Module 2", 2)
	mustExec(t, db, `INSERT INTO topics (id, title, status, tags) VALUES (?, ?, ?, ?)`,
		"rust-basics", "Rust Basics", "published", `["rust"]`)
	mustExec(t, db, `INSERT INTO concepts (id, name, definition, defined_in_topic, status) VALUES (?, ?, ?, ?, 'active')`,
		"goroutine", "Goroutine", "A lightweight thread in Go", "go-basics")
	mustExec(t, db, `INSERT INTO concepts (id, name, definition, defined_in_topic, status) VALUES (?, ?, ?, ?, 'active')`,
		"ownership", "Ownership", "Rust memory model", "rust-basics")

	builder := research.NewPoolSummaryBuilder(db)

	data, err := builder.Build(context.Background())
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	var summary research.PoolSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(summary.ExistingTopics) != 2 {
		t.Fatalf("expected 2 topics, got %d", len(summary.ExistingTopics))
	}

	// Topics ordered by ID.
	goTopic := summary.ExistingTopics[0]
	if goTopic.ID != "go-basics" {
		t.Fatalf("expected first topic 'go-basics', got %q", goTopic.ID)
	}

	if len(goTopic.Modules) != 2 {
		t.Fatalf("expected 2 modules for go-basics, got %d", len(goTopic.Modules))
	}

	// Modules ordered by sort_order.
	if goTopic.Modules[0] != "mod-1" || goTopic.Modules[1] != "mod-2" {
		t.Fatalf("expected modules [mod-1, mod-2], got %v", goTopic.Modules)
	}

	// Rust topic has no modules.
	rustTopic := summary.ExistingTopics[1]
	if len(rustTopic.Modules) != 0 {
		t.Fatalf("expected 0 modules for rust-basics, got %d", len(rustTopic.Modules))
	}

	if len(summary.ExistingConcepts) != 2 {
		t.Fatalf("expected 2 concepts, got %d", len(summary.ExistingConcepts))
	}
}

func TestPoolSummaryWriteToDir(t *testing.T) {
	db := setupTestDB(t)
	builder := research.NewPoolSummaryBuilder(db)

	dir := t.TempDir()

	if err := builder.WriteToDir(context.Background(), dir); err != nil {
		t.Fatalf("write to dir: %v", err)
	}

	path := filepath.Join(dir, "knowledge_pool_summary.json")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	var summary research.PoolSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(summary.ExistingTopics) != 0 {
		t.Fatalf("expected 0 topics, got %d", len(summary.ExistingTopics))
	}
}
