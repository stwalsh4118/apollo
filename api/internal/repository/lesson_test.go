package repository_test

import (
	"context"
	"testing"

	"github.com/sean/apollo/api/internal/repository"
)

func TestGetLessonByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewLessonRepository(db)

	seedTopic(t, db, "go-topic", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "go-topic", "Introduction", 1)
	seedLesson(t, db, "lesson-1", "mod-1", "Hello World", 1)

	lesson, err := repo.GetLessonByID(context.Background(), "lesson-1")
	if err != nil {
		t.Fatalf("get lesson: %v", err)
	}

	if lesson == nil {
		t.Fatal("expected non-nil lesson")
	}

	if lesson.Title != "Hello World" {
		t.Fatalf("expected title 'Hello World', got %q", lesson.Title)
	}

	if lesson.ModuleID != "mod-1" {
		t.Fatalf("expected module_id 'mod-1', got %q", lesson.ModuleID)
	}

	if len(lesson.Content) == 0 {
		t.Fatal("expected non-empty content JSON")
	}
}

func TestGetLessonByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewLessonRepository(db)

	lesson, err := repo.GetLessonByID(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if lesson != nil {
		t.Fatal("expected nil lesson for nonexistent ID")
	}
}
