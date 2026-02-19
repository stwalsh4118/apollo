package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/handler"
	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
)

type mockProgressRepo struct {
	topicProgress   *models.TopicProgress
	lessonProgress  *models.LessonProgress
	progressSummary *models.ProgressSummary
	returnErr       error
}

func (m *mockProgressRepo) GetTopicProgress(_ context.Context, _ string) (*models.TopicProgress, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.topicProgress, nil
}

func (m *mockProgressRepo) UpdateLessonProgress(_ context.Context, _ string, _ models.UpdateProgressInput) (*models.LessonProgress, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.lessonProgress, nil
}

func (m *mockProgressRepo) GetProgressSummary(_ context.Context) (*models.ProgressSummary, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.progressSummary, nil
}

func TestGetTopicProgressHandler(t *testing.T) {
	tp := &models.TopicProgress{
		TopicID: "topic-1",
		Lessons: []models.LessonProgress{
			{LessonID: "l1", Status: "not_started"},
		},
	}

	r := chi.NewRouter()
	handler.NewProgressHandler(&mockProgressRepo{topicProgress: tp}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/progress/topics/topic-1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result models.TopicProgress
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result.TopicID != "topic-1" {
		t.Fatalf("expected topic_id 'topic-1', got %q", result.TopicID)
	}
}

func TestGetTopicProgressNotFoundHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewProgressHandler(&mockProgressRepo{returnErr: repository.ErrNotFound}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/progress/topics/nonexistent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestUpdateLessonProgressHandler(t *testing.T) {
	lp := &models.LessonProgress{
		LessonID: "l1",
		Status:   "completed",
	}

	r := chi.NewRouter()
	handler.NewProgressHandler(&mockProgressRepo{lessonProgress: lp}).RegisterRoutes(r)

	body := `{"status":"completed","notes":"Great!"}`
	req := httptest.NewRequest(http.MethodPut, "/api/progress/lessons/l1", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var result models.LessonProgress
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result.Status != "completed" {
		t.Fatalf("expected status 'completed', got %q", result.Status)
	}
}

func TestUpdateLessonProgressBadJSON(t *testing.T) {
	r := chi.NewRouter()
	handler.NewProgressHandler(&mockProgressRepo{}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPut, "/api/progress/lessons/l1", strings.NewReader("{bad"))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUpdateLessonProgressNotFoundHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewProgressHandler(&mockProgressRepo{returnErr: repository.ErrNotFound}).RegisterRoutes(r)

	body := `{"status":"completed"}`
	req := httptest.NewRequest(http.MethodPut, "/api/progress/lessons/nonexistent", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetProgressSummaryHandler(t *testing.T) {
	ps := &models.ProgressSummary{
		TotalLessons:         10,
		CompletedLessons:     3,
		CompletionPercentage: 30.0,
		ActiveTopics:         2,
	}

	r := chi.NewRouter()
	handler.NewProgressHandler(&mockProgressRepo{progressSummary: ps}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/progress/summary", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result models.ProgressSummary
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result.TotalLessons != 10 {
		t.Fatalf("expected 10 total lessons, got %d", result.TotalLessons)
	}
}
