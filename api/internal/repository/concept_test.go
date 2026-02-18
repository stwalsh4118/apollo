package repository_test

import (
	"context"
	"testing"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
)

func TestListConcepts(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewConceptRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedConcept(t, db, "c1", "Variables", "Storage locations", "t1")
	seedConcept(t, db, "c2", "Functions", "Reusable code blocks", "t1")
	seedConcept(t, db, "c3", "Goroutines", "Lightweight threads", "t1")
	seedConcept(t, db, "c4", "Channels", "Communication pipes", "t1")
	seedConcept(t, db, "c5", "Interfaces", "Method sets", "t1")

	params := models.PaginationParams{Page: 1, PerPage: 3}

	result, err := repo.ListConcepts(context.Background(), params, "")
	if err != nil {
		t.Fatalf("list concepts: %v", err)
	}

	if result.Total != 5 {
		t.Fatalf("expected total 5, got %d", result.Total)
	}

	if len(result.Items) != 3 {
		t.Fatalf("expected 3 items on page 1, got %d", len(result.Items))
	}

	if result.Page != 1 {
		t.Fatalf("expected page 1, got %d", result.Page)
	}

	if result.PerPage != 3 {
		t.Fatalf("expected per_page 3, got %d", result.PerPage)
	}
}

func TestListConceptsPage2(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewConceptRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedConcept(t, db, "c1", "Variables", "Storage locations", "t1")
	seedConcept(t, db, "c2", "Functions", "Reusable code blocks", "t1")
	seedConcept(t, db, "c3", "Goroutines", "Lightweight threads", "t1")
	seedConcept(t, db, "c4", "Channels", "Communication pipes", "t1")
	seedConcept(t, db, "c5", "Interfaces", "Method sets", "t1")

	params := models.PaginationParams{Page: 2, PerPage: 3}

	result, err := repo.ListConcepts(context.Background(), params, "")
	if err != nil {
		t.Fatalf("list concepts page 2: %v", err)
	}

	if result.Total != 5 {
		t.Fatalf("expected total 5, got %d", result.Total)
	}

	if len(result.Items) != 2 {
		t.Fatalf("expected 2 items on page 2, got %d", len(result.Items))
	}
}

func TestListConceptsWithTopicFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewConceptRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedTopic(t, db, "t2", "Advanced Go", "advanced", "published")
	seedConcept(t, db, "c1", "Variables", "Storage locations", "t1")
	seedConcept(t, db, "c2", "Functions", "Reusable code blocks", "t1")
	seedConcept(t, db, "c3", "Goroutines", "Lightweight threads", "t2")

	params := models.PaginationParams{Page: 1, PerPage: 20}

	result, err := repo.ListConcepts(context.Background(), params, "t1")
	if err != nil {
		t.Fatalf("list concepts with filter: %v", err)
	}

	if result.Total != 2 {
		t.Fatalf("expected total 2 for topic t1, got %d", result.Total)
	}

	if len(result.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result.Items))
	}
}

func TestGetConceptByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewConceptRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedModule(t, db, "m1", "t1", "Intro", 1)
	seedLesson(t, db, "l1", "m1", "Hello World", 1)
	seedLesson(t, db, "l2", "m1", "Variables", 2)
	seedConcept(t, db, "c1", "Variables", "Storage locations", "t1")
	seedConceptReference(t, db, "c1", "l1")
	seedConceptReference(t, db, "c1", "l2")

	concept, err := repo.GetConceptByID(context.Background(), "c1")
	if err != nil {
		t.Fatalf("get concept: %v", err)
	}

	if concept == nil {
		t.Fatal("expected non-nil concept")
	}

	if concept.Name != "Variables" {
		t.Fatalf("expected name 'Variables', got %q", concept.Name)
	}

	if concept.Definition != "Storage locations" {
		t.Fatalf("expected definition 'Storage locations', got %q", concept.Definition)
	}

	if len(concept.References) != 2 {
		t.Fatalf("expected 2 references, got %d", len(concept.References))
	}
}

func TestGetConceptByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewConceptRepository(db)

	concept, err := repo.GetConceptByID(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if concept != nil {
		t.Fatal("expected nil concept for nonexistent ID")
	}
}

func TestGetConceptReferences(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewConceptRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedModule(t, db, "m1", "t1", "Intro", 1)
	seedLesson(t, db, "l1", "m1", "Hello World", 1)
	seedLesson(t, db, "l2", "m1", "Variables", 2)
	seedLesson(t, db, "l3", "m1", "Functions", 3)
	seedConcept(t, db, "c1", "Variables", "Storage locations", "t1")
	seedConceptReference(t, db, "c1", "l1")
	seedConceptReference(t, db, "c1", "l2")
	seedConceptReference(t, db, "c1", "l3")

	refs, err := repo.GetConceptReferences(context.Background(), "c1")
	if err != nil {
		t.Fatalf("get references: %v", err)
	}

	if len(refs) != 3 {
		t.Fatalf("expected 3 references, got %d", len(refs))
	}

	// Verify lesson titles are populated from join.
	for _, ref := range refs {
		if ref.LessonTitle == "" {
			t.Fatal("expected non-empty lesson title")
		}

		if ref.LessonID == "" {
			t.Fatal("expected non-empty lesson ID")
		}
	}
}

func TestGetConceptReferencesNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewConceptRepository(db)

	refs, err := repo.GetConceptReferences(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if refs != nil {
		t.Fatal("expected nil refs for nonexistent concept")
	}
}

func TestGetConceptReferencesEmpty(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewConceptRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedConcept(t, db, "c1", "Variables", "Storage locations", "t1")

	// Concept exists but has no references.
	refs, err := repo.GetConceptReferences(context.Background(), "c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if refs == nil {
		t.Fatal("expected non-nil empty slice for existing concept with no refs")
	}

	if len(refs) != 0 {
		t.Fatalf("expected 0 refs, got %d", len(refs))
	}
}
