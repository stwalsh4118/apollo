package repository_test

import (
	"context"
	"testing"

	"github.com/sean/apollo/api/internal/repository"
)

func TestGetModuleByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewModuleRepository(db)

	seedTopic(t, db, "go-topic", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "go-topic", "Introduction", 1)
	seedLesson(t, db, "lesson-1", "mod-1", "Hello World", 1)
	seedLesson(t, db, "lesson-2", "mod-1", "Variables", 2)
	seedLesson(t, db, "lesson-3", "mod-1", "Functions", 3)

	module, err := repo.GetModuleByID(context.Background(), "mod-1")
	if err != nil {
		t.Fatalf("get module: %v", err)
	}

	if module == nil {
		t.Fatal("expected non-nil module")
	}

	if module.Title != "Introduction" {
		t.Fatalf("expected title 'Introduction', got %q", module.Title)
	}

	if len(module.Lessons) != 3 {
		t.Fatalf("expected 3 lessons, got %d", len(module.Lessons))
	}

	if module.Lessons[0].Title != "Hello World" {
		t.Fatalf("expected first lesson 'Hello World', got %q", module.Lessons[0].Title)
	}

	if module.Lessons[2].Title != "Functions" {
		t.Fatalf("expected third lesson 'Functions', got %q", module.Lessons[2].Title)
	}
}

func TestGetModuleByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewModuleRepository(db)

	module, err := repo.GetModuleByID(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if module != nil {
		t.Fatal("expected nil module for nonexistent ID")
	}
}
