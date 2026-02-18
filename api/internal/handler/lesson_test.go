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

type mockLessonRepo struct {
	detail    *models.LessonDetail
	returnErr error
}

func (m *mockLessonRepo) GetLessonByID(_ context.Context, _ string) (*models.LessonDetail, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.detail, nil
}

func TestGetLessonByIDHandler(t *testing.T) {
	ld := &models.LessonDetail{}
	ld.ID = "l1"
	ld.ModuleID = "m1"
	ld.Title = "Lesson 1"
	ld.Content = []byte(`{"text":"hello"}`)

	r := chi.NewRouter()
	handler.NewLessonHandler(&mockLessonRepo{detail: ld}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/lessons/l1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result["id"] != "l1" {
		t.Fatalf("expected id 'l1', got %v", result["id"])
	}
}

func TestGetLessonByIDNotFoundHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewLessonHandler(&mockLessonRepo{detail: nil}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/lessons/nonexistent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetLessonByIDErrorHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewLessonHandler(&mockLessonRepo{returnErr: errors.New("db error")}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/lessons/l1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
