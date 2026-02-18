package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
)

func TestSearch(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSearchRepository(db)

	mustExec(t, db, "INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)",
		"topic", "t1", "Go Basics", "Learn the fundamentals of Go programming")
	mustExec(t, db, "INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)",
		"lesson", "l1", "Hello World", "Write your first Go program")
	mustExec(t, db, "INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)",
		"concept", "c1", "Variable", "A named storage location in Go")

	params := models.PaginationParams{Page: 1, PerPage: 20}

	result, err := repo.Search(context.Background(), "Go", params)
	if err != nil {
		t.Fatalf("search: %v", err)
	}

	if result.Total < 2 {
		t.Fatalf("expected at least 2 results for 'Go', got %d", result.Total)
	}

	if len(result.Items) < 2 {
		t.Fatalf("expected at least 2 items, got %d", len(result.Items))
	}

	for _, item := range result.Items {
		if item.EntityType == "" {
			t.Fatal("expected non-empty entity_type")
		}

		if item.EntityID == "" {
			t.Fatal("expected non-empty entity_id")
		}

		if item.Title == "" {
			t.Fatal("expected non-empty title")
		}
	}
}

func TestSearchPagination(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSearchRepository(db)

	for i := 0; i < 10; i++ {
		mustExec(t, db,
			"INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)",
			"topic", fmt.Sprintf("t%d", i), fmt.Sprintf("Go Topic %d", i), "Go programming fundamentals",
		)
	}

	params := models.PaginationParams{Page: 1, PerPage: 3}

	result, err := repo.Search(context.Background(), "Go", params)
	if err != nil {
		t.Fatalf("search page 1: %v", err)
	}

	if result.Total != 10 {
		t.Fatalf("expected total 10, got %d", result.Total)
	}

	if len(result.Items) != 3 {
		t.Fatalf("expected 3 items on page 1, got %d", len(result.Items))
	}

	params2 := models.PaginationParams{Page: 4, PerPage: 3}

	result2, err := repo.Search(context.Background(), "Go", params2)
	if err != nil {
		t.Fatalf("search page 4: %v", err)
	}

	if len(result2.Items) != 1 {
		t.Fatalf("expected 1 item on page 4, got %d", len(result2.Items))
	}
}

func TestSearchNoResults(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSearchRepository(db)

	mustExec(t, db, "INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)",
		"topic", "t1", "Go Basics", "Learn Go programming")

	params := models.PaginationParams{Page: 1, PerPage: 20}

	result, err := repo.Search(context.Background(), "xyznonexistent", params)
	if err != nil {
		t.Fatalf("search: %v", err)
	}

	if result.Total != 0 {
		t.Fatalf("expected total 0, got %d", result.Total)
	}

	if len(result.Items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(result.Items))
	}
}

func TestSearchInvalidFTS5Syntax(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSearchRepository(db)

	mustExec(t, db, "INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)",
		"topic", "t1", "Go Basics", "Learn Go programming")

	params := models.PaginationParams{Page: 1, PerPage: 20}

	_, err := repo.Search(context.Background(), "\"unclosed", params)
	if err == nil {
		t.Fatal("expected error for invalid FTS5 syntax")
	}

	if !errors.Is(err, repository.ErrInvalidQuery) {
		t.Fatalf("expected ErrInvalidQuery, got: %v", err)
	}
}

func TestSearchSnippet(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSearchRepository(db)

	mustExec(t, db, "INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)",
		"topic", "t1", "Go Basics", "Learn the fundamentals of Go programming language")

	params := models.PaginationParams{Page: 1, PerPage: 20}

	result, err := repo.Search(context.Background(), "fundamentals", params)
	if err != nil {
		t.Fatalf("search: %v", err)
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Items))
	}

	if result.Items[0].Snippet == "" {
		t.Fatal("expected non-empty snippet")
	}
}
