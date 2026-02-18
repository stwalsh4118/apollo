package repository_test

import (
	"context"
	"testing"

	"github.com/sean/apollo/api/internal/repository"
)

func TestListTopicsEmpty(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTopicRepository(db)

	topics, err := repo.ListTopics(context.Background())
	if err != nil {
		t.Fatalf("list topics: %v", err)
	}

	if len(topics) != 0 {
		t.Fatalf("expected 0 topics, got %d", len(topics))
	}
}

func TestListTopicsOrderedByTitle(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTopicRepository(db)

	seedTopic(t, db, "c-topic", "C Programming", "foundational", "published")
	seedTopic(t, db, "a-topic", "Algorithms", "intermediate", "published")
	seedTopic(t, db, "b-topic", "Bash Scripting", "foundational", "published")

	topics, err := repo.ListTopics(context.Background())
	if err != nil {
		t.Fatalf("list topics: %v", err)
	}

	if len(topics) != 3 {
		t.Fatalf("expected 3 topics, got %d", len(topics))
	}

	if topics[0].Title != "Algorithms" {
		t.Fatalf("expected first topic 'Algorithms', got %q", topics[0].Title)
	}

	if topics[1].Title != "Bash Scripting" {
		t.Fatalf("expected second topic 'Bash Scripting', got %q", topics[1].Title)
	}
}

func TestListTopicsWithModuleCount(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTopicRepository(db)

	seedTopic(t, db, "go-topic", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "go-topic", "Intro", 1)
	seedModule(t, db, "mod-2", "go-topic", "Types", 2)

	topics, err := repo.ListTopics(context.Background())
	if err != nil {
		t.Fatalf("list topics: %v", err)
	}

	if topics[0].ModuleCount != 2 {
		t.Fatalf("expected module_count 2, got %d", topics[0].ModuleCount)
	}
}

func TestListTopicsJSONTags(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTopicRepository(db)

	seedTopic(t, db, "go-topic", "Go Basics", "foundational", "published")

	topics, err := repo.ListTopics(context.Background())
	if err != nil {
		t.Fatalf("list topics: %v", err)
	}

	if len(topics[0].Tags) != 1 || topics[0].Tags[0] != "test" {
		t.Fatalf("expected tags [test], got %v", topics[0].Tags)
	}
}

func TestGetTopicByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTopicRepository(db)

	seedTopic(t, db, "go-topic", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "go-topic", "Introduction", 1)
	seedModule(t, db, "mod-2", "go-topic", "Types", 2)

	topic, err := repo.GetTopicByID(context.Background(), "go-topic")
	if err != nil {
		t.Fatalf("get topic: %v", err)
	}

	if topic == nil {
		t.Fatal("expected non-nil topic")
	}

	if topic.Title != "Go Basics" {
		t.Fatalf("expected title 'Go Basics', got %q", topic.Title)
	}

	if len(topic.Modules) != 2 {
		t.Fatalf("expected 2 modules, got %d", len(topic.Modules))
	}

	if topic.Modules[0].Title != "Introduction" {
		t.Fatalf("expected first module 'Introduction', got %q", topic.Modules[0].Title)
	}
}

func TestGetTopicByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTopicRepository(db)

	topic, err := repo.GetTopicByID(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if topic != nil {
		t.Fatal("expected nil topic for nonexistent ID")
	}
}

func TestGetTopicFull(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTopicRepository(db)

	seedTopic(t, db, "go-topic", "Go Basics", "foundational", "published")
	seedModule(t, db, "mod-1", "go-topic", "Introduction", 1)
	seedLesson(t, db, "lesson-1", "mod-1", "Hello World", 1)
	seedConcept(t, db, "concept-1", "Variable", "A named storage", "go-topic")
	seedConceptReference(t, db, "concept-1", "lesson-1")

	topic, err := repo.GetTopicFull(context.Background(), "go-topic")
	if err != nil {
		t.Fatalf("get topic full: %v", err)
	}

	if topic == nil {
		t.Fatal("expected non-nil topic")
	}

	if len(topic.Modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(topic.Modules))
	}

	mod := topic.Modules[0]
	if len(mod.Lessons) != 1 {
		t.Fatalf("expected 1 lesson, got %d", len(mod.Lessons))
	}

	lesson := mod.Lessons[0]
	if lesson.Title != "Hello World" {
		t.Fatalf("expected lesson title 'Hello World', got %q", lesson.Title)
	}

	if len(lesson.Content) == 0 {
		t.Fatal("expected non-empty content JSON")
	}

	if len(lesson.Concepts) != 1 {
		t.Fatalf("expected 1 concept, got %d", len(lesson.Concepts))
	}

	if lesson.Concepts[0].Name != "Variable" {
		t.Fatalf("expected concept 'Variable', got %q", lesson.Concepts[0].Name)
	}
}

func TestGetTopicFullNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTopicRepository(db)

	topic, err := repo.GetTopicFull(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if topic != nil {
		t.Fatal("expected nil topic for nonexistent ID")
	}
}
