package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
)

func TestCreateTopic(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)
	readRepo := repository.NewTopicRepository(db)

	input := models.TopicInput{
		ID:         "go-basics",
		Title:      "Go Basics",
		Difficulty: "foundational",
		Status:     "published",
		Tags:       []string{"go", "basics"},
	}

	if err := writeRepo.CreateTopic(context.Background(), input); err != nil {
		t.Fatalf("create topic: %v", err)
	}

	topics, err := readRepo.ListTopics(context.Background())
	if err != nil {
		t.Fatalf("list topics: %v", err)
	}

	if len(topics) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(topics))
	}

	if topics[0].Title != "Go Basics" {
		t.Fatalf("expected title 'Go Basics', got %q", topics[0].Title)
	}
}

func TestUpdateTopic(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)
	readRepo := repository.NewTopicRepository(db)

	input := models.TopicInput{
		ID:     "go-basics",
		Title:  "Go Basics",
		Status: "draft",
	}

	if err := writeRepo.CreateTopic(context.Background(), input); err != nil {
		t.Fatalf("create topic: %v", err)
	}

	updated := models.TopicInput{
		Title:  "Go Fundamentals",
		Status: "published",
	}

	if err := writeRepo.UpdateTopic(context.Background(), "go-basics", updated); err != nil {
		t.Fatalf("update topic: %v", err)
	}

	topics, err := readRepo.ListTopics(context.Background())
	if err != nil {
		t.Fatalf("list topics: %v", err)
	}

	if topics[0].Title != "Go Fundamentals" {
		t.Fatalf("expected title 'Go Fundamentals', got %q", topics[0].Title)
	}
}

func TestUpdateTopicNotFound(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)

	err := writeRepo.UpdateTopic(context.Background(), "nonexistent", models.TopicInput{Title: "X", Status: "draft"})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateModule(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)
	readRepo := repository.NewModuleRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")

	input := models.ModuleInput{
		ID:        "m1",
		TopicID:   "t1",
		Title:     "Introduction",
		SortOrder: 1,
	}

	if err := writeRepo.CreateModule(context.Background(), input); err != nil {
		t.Fatalf("create module: %v", err)
	}

	// Seed a lesson so GetModuleByID has something to return.
	seedLesson(t, db, "l1", "m1", "Hello World", 1)

	module, err := readRepo.GetModuleByID(context.Background(), "m1")
	if err != nil {
		t.Fatalf("get module: %v", err)
	}

	if module.Title != "Introduction" {
		t.Fatalf("expected title 'Introduction', got %q", module.Title)
	}
}

func TestCreateLesson(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)
	readRepo := repository.NewLessonRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedModule(t, db, "m1", "t1", "Intro", 1)

	input := models.LessonInput{
		ID:       "l1",
		ModuleID: "m1",
		Title:    "Hello World",
		Content:  []byte(`[{"type":"text","body":"Hello"}]`),
	}

	if err := writeRepo.CreateLesson(context.Background(), input); err != nil {
		t.Fatalf("create lesson: %v", err)
	}

	lesson, err := readRepo.GetLessonByID(context.Background(), "l1")
	if err != nil {
		t.Fatalf("get lesson: %v", err)
	}

	if lesson.Title != "Hello World" {
		t.Fatalf("expected title 'Hello World', got %q", lesson.Title)
	}
}

func TestCreateConcept(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)
	readRepo := repository.NewConceptRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")

	input := models.ConceptInput{
		ID:             "c1",
		Name:           "Variable",
		Definition:     "A named storage location",
		DefinedInTopic: "t1",
		Aliases:        []string{"var"},
	}

	if err := writeRepo.CreateConcept(context.Background(), input); err != nil {
		t.Fatalf("create concept: %v", err)
	}

	concept, err := readRepo.GetConceptByID(context.Background(), "c1")
	if err != nil {
		t.Fatalf("get concept: %v", err)
	}

	if concept.Name != "Variable" {
		t.Fatalf("expected name 'Variable', got %q", concept.Name)
	}
}

func TestCreateConceptReference(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)
	readRepo := repository.NewConceptRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedModule(t, db, "m1", "t1", "Intro", 1)
	seedLesson(t, db, "l1", "m1", "Hello World", 1)
	seedConcept(t, db, "c1", "Variable", "A named storage location", "t1")

	input := models.ConceptReferenceInput{
		LessonID: "l1",
		Context:  "Used in examples",
	}

	if err := writeRepo.CreateConceptReference(context.Background(), "c1", input); err != nil {
		t.Fatalf("create concept reference: %v", err)
	}

	refs, err := readRepo.GetConceptReferences(context.Background(), "c1")
	if err != nil {
		t.Fatalf("get references: %v", err)
	}

	if len(refs) != 1 {
		t.Fatalf("expected 1 reference, got %d", len(refs))
	}

	if refs[0].LessonTitle != "Hello World" {
		t.Fatalf("expected lesson title 'Hello World', got %q", refs[0].LessonTitle)
	}
}

func TestCreatePrerequisite(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedTopic(t, db, "t2", "Advanced Go", "advanced", "published")

	input := models.PrerequisiteInput{
		TopicID:             "t2",
		PrerequisiteTopicID: "t1",
		Priority:            "essential",
		Reason:              "Must know basics first",
	}

	if err := writeRepo.CreatePrerequisite(context.Background(), input); err != nil {
		t.Fatalf("create prerequisite: %v", err)
	}
}

func TestCreateRelation(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedTopic(t, db, "t2", "Advanced Go", "advanced", "published")

	input := models.RelationInput{
		TopicA:       "t1",
		TopicB:       "t2",
		RelationType: "builds_on",
		Description:  "Advanced topics build on basics",
	}

	if err := writeRepo.CreateRelation(context.Background(), input); err != nil {
		t.Fatalf("create relation: %v", err)
	}
}

func TestCreateTopicDuplicate(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)

	input := models.TopicInput{
		ID:     "t1",
		Title:  "Go Basics",
		Status: "draft",
	}

	if err := writeRepo.CreateTopic(context.Background(), input); err != nil {
		t.Fatalf("first create: %v", err)
	}

	err := writeRepo.CreateTopic(context.Background(), input)
	if !errors.Is(err, repository.ErrDuplicate) {
		t.Fatalf("expected ErrDuplicate, got %v", err)
	}
}

func TestCreateModuleFKViolation(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)

	input := models.ModuleInput{
		ID:        "m1",
		TopicID:   "nonexistent",
		Title:     "Intro",
		SortOrder: 1,
	}

	err := writeRepo.CreateModule(context.Background(), input)
	if !errors.Is(err, repository.ErrFKViolation) {
		t.Fatalf("expected ErrFKViolation, got %v", err)
	}
}

func TestCreateTopicSearchIndex(t *testing.T) {
	db := setupTestDB(t)
	writeRepo := repository.NewWriteRepository(db)

	input := models.TopicInput{
		ID:          "t1",
		Title:       "Go Basics",
		Description: "Learn the fundamentals of Go",
		Status:      "published",
	}

	if err := writeRepo.CreateTopic(context.Background(), input); err != nil {
		t.Fatalf("create topic: %v", err)
	}

	// Verify search index entry was created.
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM search_index WHERE entity_type = 'topic' AND entity_id = 't1'").Scan(&count)
	if err != nil {
		t.Fatalf("query search index: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 search index entry, got %d", count)
	}
}
