package repository_test

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
)

func TestGetTopicProgress_DefaultNotStarted(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProgressRepository(db)

	seedTopic(t, db, "topic-1", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "topic-1", "Module 1", 1)
	seedLesson(t, db, "lesson-1", "mod-1", "Lesson 1", 1)
	seedLesson(t, db, "lesson-2", "mod-1", "Lesson 2", 2)

	tp, err := repo.GetTopicProgress(context.Background(), "topic-1")
	if err != nil {
		t.Fatalf("get topic progress: %v", err)
	}

	if tp.TopicID != "topic-1" {
		t.Fatalf("expected topic_id 'topic-1', got %q", tp.TopicID)
	}

	if len(tp.Lessons) != 2 {
		t.Fatalf("expected 2 lessons, got %d", len(tp.Lessons))
	}

	for _, lp := range tp.Lessons {
		if lp.Status != models.ProgressStatusNotStarted {
			t.Fatalf("expected status 'not_started', got %q for lesson %s", lp.Status, lp.LessonID)
		}
	}
}

func TestGetTopicProgress_NotFoundTopic(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProgressRepository(db)

	_, err := repo.GetTopicProgress(context.Background(), "nonexistent")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateLessonProgress_CreateNew(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProgressRepository(db)

	seedTopic(t, db, "topic-1", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "topic-1", "Module 1", 1)
	seedLesson(t, db, "lesson-1", "mod-1", "Lesson 1", 1)

	lp, err := repo.UpdateLessonProgress(context.Background(), "lesson-1", models.UpdateProgressInput{
		Status: models.ProgressStatusInProgress,
	})
	if err != nil {
		t.Fatalf("update lesson progress: %v", err)
	}

	if lp.LessonID != "lesson-1" {
		t.Fatalf("expected lesson_id 'lesson-1', got %q", lp.LessonID)
	}

	if lp.Status != models.ProgressStatusInProgress {
		t.Fatalf("expected status 'in_progress', got %q", lp.Status)
	}

	if lp.StartedAt == "" {
		t.Fatal("expected started_at to be set")
	}

	if lp.CompletedAt != "" {
		t.Fatalf("expected empty completed_at, got %q", lp.CompletedAt)
	}
}

func TestUpdateLessonProgress_Upsert(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProgressRepository(db)

	seedTopic(t, db, "topic-1", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "topic-1", "Module 1", 1)
	seedLesson(t, db, "lesson-1", "mod-1", "Lesson 1", 1)

	// First update: in_progress.
	_, err := repo.UpdateLessonProgress(context.Background(), "lesson-1", models.UpdateProgressInput{
		Status: models.ProgressStatusInProgress,
	})
	if err != nil {
		t.Fatalf("first update: %v", err)
	}

	// Second update: completed with notes.
	lp, err := repo.UpdateLessonProgress(context.Background(), "lesson-1", models.UpdateProgressInput{
		Status: models.ProgressStatusCompleted,
		Notes:  "Great lesson!",
	})
	if err != nil {
		t.Fatalf("second update: %v", err)
	}

	if lp.Status != models.ProgressStatusCompleted {
		t.Fatalf("expected status 'completed', got %q", lp.Status)
	}

	if lp.Notes != "Great lesson!" {
		t.Fatalf("expected notes 'Great lesson!', got %q", lp.Notes)
	}

	if lp.StartedAt == "" {
		t.Fatal("expected started_at to be preserved from first update")
	}

	if lp.CompletedAt == "" {
		t.Fatal("expected completed_at to be set")
	}
}

func TestUpdateLessonProgress_NotFoundLesson(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProgressRepository(db)

	_, err := repo.UpdateLessonProgress(context.Background(), "nonexistent", models.UpdateProgressInput{
		Status: models.ProgressStatusCompleted,
	})

	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetProgressSummary(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProgressRepository(db)

	seedTopic(t, db, "topic-1", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "topic-1", "Module 1", 1)
	seedLesson(t, db, "lesson-1", "mod-1", "Lesson 1", 1)
	seedLesson(t, db, "lesson-2", "mod-1", "Lesson 2", 2)
	seedLesson(t, db, "lesson-3", "mod-1", "Lesson 3", 3)

	// Complete one lesson.
	_, err := repo.UpdateLessonProgress(context.Background(), "lesson-1", models.UpdateProgressInput{
		Status: models.ProgressStatusCompleted,
	})
	if err != nil {
		t.Fatalf("update progress: %v", err)
	}

	ps, err := repo.GetProgressSummary(context.Background())
	if err != nil {
		t.Fatalf("get progress summary: %v", err)
	}

	if ps.TotalLessons != 3 {
		t.Fatalf("expected 3 total lessons, got %d", ps.TotalLessons)
	}

	if ps.CompletedLessons != 1 {
		t.Fatalf("expected 1 completed lesson, got %d", ps.CompletedLessons)
	}

	expectedPct := 100.0 / 3.0
	if math.Abs(ps.CompletionPercentage-expectedPct) > 0.01 {
		t.Fatalf("expected completion ~%.2f%%, got %.2f%%", expectedPct, ps.CompletionPercentage)
	}

	if ps.ActiveTopics != 1 {
		t.Fatalf("expected 1 active topic, got %d", ps.ActiveTopics)
	}
}

func TestCompletedAtAutoSet(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProgressRepository(db)

	seedTopic(t, db, "topic-1", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "topic-1", "Module 1", 1)
	seedLesson(t, db, "lesson-1", "mod-1", "Lesson 1", 1)

	lp, err := repo.UpdateLessonProgress(context.Background(), "lesson-1", models.UpdateProgressInput{
		Status: models.ProgressStatusCompleted,
	})
	if err != nil {
		t.Fatalf("update progress: %v", err)
	}

	if lp.CompletedAt == "" {
		t.Fatal("expected completed_at to be auto-set on completed status")
	}

	if lp.StartedAt == "" {
		t.Fatal("expected started_at to be auto-set on completed status")
	}
}
