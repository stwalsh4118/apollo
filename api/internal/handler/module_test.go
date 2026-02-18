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

type mockModuleRepo struct {
	detail    *models.ModuleDetail
	returnErr error
}

func (m *mockModuleRepo) GetModuleByID(_ context.Context, _ string) (*models.ModuleDetail, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.detail, nil
}

func TestGetModuleByIDHandler(t *testing.T) {
	md := &models.ModuleDetail{
		Lessons: []models.LessonSummary{{ID: "l1", Title: "Lesson 1", SortOrder: 1}},
	}
	md.ID = "m1"
	md.TopicID = "t1"
	md.Title = "Module 1"

	r := chi.NewRouter()
	handler.NewModuleHandler(&mockModuleRepo{detail: md}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/modules/m1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result["id"] != "m1" {
		t.Fatalf("expected id 'm1', got %v", result["id"])
	}
}

func TestGetModuleByIDNotFoundHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewModuleHandler(&mockModuleRepo{detail: nil}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/modules/nonexistent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetModuleByIDErrorHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewModuleHandler(&mockModuleRepo{returnErr: errors.New("db error")}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/modules/m1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
