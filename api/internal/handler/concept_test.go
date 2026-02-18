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

type mockConceptRepo struct {
	listResult *models.PaginatedResponse[models.ConceptSummary]
	detail     *models.ConceptDetail
	refs       []models.ConceptReference
	returnErr  error
}

func (m *mockConceptRepo) ListConcepts(_ context.Context, _ models.PaginationParams, _ string) (*models.PaginatedResponse[models.ConceptSummary], error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.listResult, nil
}

func (m *mockConceptRepo) GetConceptByID(_ context.Context, _ string) (*models.ConceptDetail, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.detail, nil
}

func (m *mockConceptRepo) GetConceptReferences(_ context.Context, _ string) ([]models.ConceptReference, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.refs, nil
}

func TestListConceptsHandler(t *testing.T) {
	mock := &mockConceptRepo{
		listResult: &models.PaginatedResponse[models.ConceptSummary]{
			Items:   []models.ConceptSummary{{ID: "c1", Name: "Variables", Status: "active"}},
			Total:   1,
			Page:    1,
			PerPage: 20,
		},
	}

	r := chi.NewRouter()
	handler.NewConceptHandler(mock).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts", nil)
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
		t.Fatal("expected items array in response")
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestListConceptsWithTopicFilterHandler(t *testing.T) {
	mock := &mockConceptRepo{
		listResult: &models.PaginatedResponse[models.ConceptSummary]{
			Items:   []models.ConceptSummary{},
			Total:   0,
			Page:    1,
			PerPage: 20,
		},
	}

	r := chi.NewRouter()
	handler.NewConceptHandler(mock).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts?topic=t1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestListConceptsPaginationHandler(t *testing.T) {
	mock := &mockConceptRepo{
		listResult: &models.PaginatedResponse[models.ConceptSummary]{
			Items:   []models.ConceptSummary{},
			Total:   5,
			Page:    2,
			PerPage: 2,
		},
	}

	r := chi.NewRouter()
	handler.NewConceptHandler(mock).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts?page=2&per_page=2", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestListConceptsErrorHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewConceptHandler(&mockConceptRepo{returnErr: errors.New("db error")}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetConceptByIDHandler(t *testing.T) {
	cd := &models.ConceptDetail{
		ID:         "c1",
		Name:       "Variables",
		Definition: "Storage locations",
		Status:     "active",
		References: []models.ConceptReference{
			{LessonID: "l1", LessonTitle: "Intro", Context: "test"},
		},
	}

	r := chi.NewRouter()
	handler.NewConceptHandler(&mockConceptRepo{detail: cd}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts/c1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result["id"] != "c1" {
		t.Fatalf("expected id 'c1', got %v", result["id"])
	}

	if result["name"] != "Variables" {
		t.Fatalf("expected name 'Variables', got %v", result["name"])
	}
}

func TestGetConceptByIDNotFoundHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewConceptHandler(&mockConceptRepo{detail: nil}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts/nonexistent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetConceptByIDErrorHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewConceptHandler(&mockConceptRepo{returnErr: errors.New("db error")}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts/c1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetConceptReferencesHandler(t *testing.T) {
	refs := []models.ConceptReference{
		{LessonID: "l1", LessonTitle: "Hello World", Context: "intro"},
		{LessonID: "l2", LessonTitle: "Variables", Context: "definition"},
	}

	r := chi.NewRouter()
	handler.NewConceptHandler(&mockConceptRepo{refs: refs}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts/c1/references", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result []any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 references, got %d", len(result))
	}
}

func TestGetConceptReferencesNotFoundHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewConceptHandler(&mockConceptRepo{refs: nil}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts/nonexistent/references", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetConceptReferencesErrorHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewConceptHandler(&mockConceptRepo{returnErr: errors.New("db error")}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/concepts/c1/references", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
