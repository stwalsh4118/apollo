package repository_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/internal/database"
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

func seedTopic(t *testing.T, db *sql.DB, id, title, difficulty, status string) {
	t.Helper()

	mustExec(t, db,
		`INSERT INTO topics (id, title, difficulty, status, tags) VALUES (?, ?, ?, ?, ?)`,
		id, title, difficulty, status, `["test"]`,
	)
}

func seedModule(t *testing.T, db *sql.DB, id, topicID, title string, sortOrder int) {
	t.Helper()

	mustExec(t, db,
		`INSERT INTO modules (id, topic_id, title, sort_order) VALUES (?, ?, ?, ?)`,
		id, topicID, title, sortOrder,
	)
}

func seedLesson(t *testing.T, db *sql.DB, id, moduleID, title string, sortOrder int) {
	t.Helper()

	mustExec(t, db,
		`INSERT INTO lessons (id, module_id, title, sort_order, content) VALUES (?, ?, ?, ?, ?)`,
		id, moduleID, title, sortOrder, `[{"type":"text","body":"Hello"}]`,
	)
}

func seedConcept(t *testing.T, db *sql.DB, id, name, definition, topicID string) {
	t.Helper()

	mustExec(t, db,
		`INSERT INTO concepts (id, name, definition, defined_in_topic, status) VALUES (?, ?, ?, ?, 'active')`,
		id, name, definition, topicID,
	)
}

func seedConceptReference(t *testing.T, db *sql.DB, conceptID, lessonID string) {
	t.Helper()

	mustExec(t, db,
		`INSERT INTO concept_references (concept_id, lesson_id, context) VALUES (?, ?, 'test context')`,
		conceptID, lessonID,
	)
}
