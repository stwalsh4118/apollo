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

// mockTopicRepo is a test double for repository.TopicRepository.
type mockTopicRepo struct {
	topics    []models.TopicSummary
	detail    *models.TopicDetail
	full      *models.TopicFull
	returnErr error
}

func (m *mockTopicRepo) ListTopics(_ context.Context) ([]models.TopicSummary, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.topics, nil
}

func (m *mockTopicRepo) GetTopicByID(_ context.Context, _ string) (*models.TopicDetail, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.detail, nil
}

func (m *mockTopicRepo) GetTopicFull(_ context.Context, _ string) (*models.TopicFull, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.full, nil
}

func setupTopicRouter(mock *mockTopicRepo) chi.Router {
	r := chi.NewRouter()
	h := handler.NewTopicHandler(mock)
	h.RegisterRoutes(r)

	return r
}

func TestListTopicsHandler(t *testing.T) {
	mock := &mockTopicRepo{
		topics: []models.TopicSummary{
			{ID: "t1", Title: "Topic 1", Status: "published", ModuleCount: 2},
		},
	}
	r := setupTopicRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/topics", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var topics []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&topics); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(topics) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(topics))
	}

	if topics[0]["id"] != "t1" {
		t.Fatalf("expected id 't1', got %v", topics[0]["id"])
	}
}

func TestListTopicsHandlerError(t *testing.T) {
	mock := &mockTopicRepo{returnErr: errors.New("db error")}
	r := setupTopicRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/topics", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetTopicByIDHandler(t *testing.T) {
	td := &models.TopicDetail{
		Modules: []models.ModuleSummary{{ID: "m1", Title: "Module 1", SortOrder: 1}},
	}
	td.ID = "t1"
	td.Title = "Topic 1"
	td.Status = "published"

	mock := &mockTopicRepo{detail: td}
	r := setupTopicRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/topics/t1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var topic map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&topic); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if topic["id"] != "t1" {
		t.Fatalf("expected id 't1', got %v", topic["id"])
	}

	modules := topic["modules"].([]any)
	if len(modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(modules))
	}
}

func TestGetTopicByIDNotFoundHandler(t *testing.T) {
	mock := &mockTopicRepo{detail: nil}
	r := setupTopicRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/topics/nonexistent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetTopicFullHandler(t *testing.T) {
	tf := &models.TopicFull{
		Modules: []models.ModuleFull{
			{
				Lessons: []models.LessonFull{
					{Concepts: []models.ConceptSummary{}},
				},
			},
		},
	}
	tf.ID = "t1"
	tf.Title = "Topic 1"
	tf.Status = "published"
	tf.Modules[0].ID = "m1"
	tf.Modules[0].Title = "Module 1"
	tf.Modules[0].Lessons[0].ID = "l1"
	tf.Modules[0].Lessons[0].Title = "Lesson 1"

	mock := &mockTopicRepo{full: tf}
	r := setupTopicRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/topics/t1/full", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var topic map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&topic); err != nil {
		t.Fatalf("decode: %v", err)
	}

	modules := topic["modules"].([]any)
	if len(modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(modules))
	}

	mod := modules[0].(map[string]any)
	lessons := mod["lessons"].([]any)
	if len(lessons) != 1 {
		t.Fatalf("expected 1 lesson, got %d", len(lessons))
	}
}

func TestGetTopicFullNotFoundHandler(t *testing.T) {
	mock := &mockTopicRepo{full: nil}
	r := setupTopicRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/topics/nonexistent/full", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
