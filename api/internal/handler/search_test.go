package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/handler"
	"github.com/sean/apollo/api/internal/models"
)

type mockSearchRepo struct {
	result    *models.PaginatedResponse[models.SearchResult]
	returnErr error
}

func (m *mockSearchRepo) Search(_ context.Context, _ string, _ models.PaginationParams) (*models.PaginatedResponse[models.SearchResult], error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.result, nil
}

func TestSearchHandler(t *testing.T) {
	mock := &mockSearchRepo{
		result: &models.PaginatedResponse[models.SearchResult]{
			Items: []models.SearchResult{
				{EntityType: "topic", EntityID: "t1", Title: "Go Basics", Snippet: "Learn <mark>Go</mark>"},
			},
			Total:   1,
			Page:    1,
			PerPage: 20,
		},
	}

	r := chi.NewRouter()
	handler.NewSearchHandler(mock).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=Go", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	items, ok := result["items"].([]any)
	if !ok {
		t.Fatal("expected items array")
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestSearchHandlerEmptyQuery(t *testing.T) {
	r := chi.NewRouter()
	handler.NewSearchHandler(&mockSearchRepo{}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestSearchHandlerError(t *testing.T) {
	r := chi.NewRouter()
	handler.NewSearchHandler(&mockSearchRepo{returnErr: errors.New("db error")}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
