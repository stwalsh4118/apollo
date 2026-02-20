package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/handler"
	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
)

type mockResearchRepo struct {
	job       *models.ResearchJob
	jobs      *models.PaginatedResponse[models.ResearchJobSummary]
	returnErr error
}

func (m *mockResearchRepo) CreateJob(_ context.Context, input models.CreateResearchJobInput) (*models.ResearchJob, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return &models.ResearchJob{
		ID:        "job-1",
		RootTopic: input.Topic,
		Status:    models.ResearchStatusQueued,
	}, nil
}

func (m *mockResearchRepo) GetJobByID(_ context.Context, _ string) (*models.ResearchJob, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.job, nil
}

func (m *mockResearchRepo) ListJobs(_ context.Context, _ models.PaginationParams) (*models.PaginatedResponse[models.ResearchJobSummary], error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.jobs, nil
}

func (m *mockResearchRepo) UpdateJobStatus(_ context.Context, _ string, status string, _ string) error {
	if m.returnErr != nil {
		return m.returnErr
	}

	if m.job != nil {
		m.job.Status = status
	}

	return nil
}

func (m *mockResearchRepo) UpdateJobProgress(_ context.Context, _ string, _ models.ResearchProgress) error {
	return m.returnErr
}

func (m *mockResearchRepo) UpdateJobCurrentTopic(_ context.Context, _ string, _ string) error {
	return m.returnErr
}

func (m *mockResearchRepo) FindOldestByStatus(_ context.Context, _ string) (string, error) {
	return "", m.returnErr
}

func TestCreateResearchJob(t *testing.T) {
	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{}, nil).RegisterRoutes(r)

	body := `{"topic":"Go Concurrency"}`
	req := httptest.NewRequest(http.MethodPost, "/api/research", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var result models.ResearchJob
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result.ID != "job-1" {
		t.Fatalf("expected job ID 'job-1', got %q", result.ID)
	}

	if result.Status != models.ResearchStatusQueued {
		t.Fatalf("expected status %q, got %q", models.ResearchStatusQueued, result.Status)
	}
}

func TestCreateResearchJobMissingTopic(t *testing.T) {
	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{}, nil).RegisterRoutes(r)

	body := `{"topic":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/research", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateResearchJobBadJSON(t *testing.T) {
	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{}, nil).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPost, "/api/research", strings.NewReader("{bad"))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestListResearchJobs(t *testing.T) {
	jobs := &models.PaginatedResponse[models.ResearchJobSummary]{
		Items:   []models.ResearchJobSummary{{ID: "j1", RootTopic: "Test", Status: "queued"}},
		Total:   1,
		Page:    1,
		PerPage: 20,
	}

	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{jobs: jobs}, nil).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/research/jobs", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result models.PaginatedResponse[models.ResearchJobSummary]
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}
}

func TestListResearchJobsEmpty(t *testing.T) {
	jobs := &models.PaginatedResponse[models.ResearchJobSummary]{
		Items:   []models.ResearchJobSummary{},
		Total:   0,
		Page:    1,
		PerPage: 20,
	}

	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{jobs: jobs}, nil).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/research/jobs", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestGetResearchJob(t *testing.T) {
	job := &models.ResearchJob{
		ID:        "job-1",
		RootTopic: "Test Topic",
		Status:    models.ResearchStatusResearching,
	}

	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{job: job}, nil).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/research/jobs/job-1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result models.ResearchJob
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result.ID != "job-1" {
		t.Fatalf("expected ID 'job-1', got %q", result.ID)
	}
}

func TestGetResearchJobNotFound(t *testing.T) {
	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{returnErr: fmt.Errorf("job: %w", repository.ErrNotFound)}, nil).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/research/jobs/nonexistent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestCancelResearchJob(t *testing.T) {
	job := &models.ResearchJob{
		ID:        "job-1",
		RootTopic: "Test",
		Status:    models.ResearchStatusResearching,
	}

	var cancelledID string
	cancelFn := func(id string) { cancelledID = id }

	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{job: job}, cancelFn).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPost, "/api/research/jobs/job-1/cancel", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	if cancelledID != "job-1" {
		t.Fatalf("expected cancel called with 'job-1', got %q", cancelledID)
	}
}

func TestCancelResearchJobNotFound(t *testing.T) {
	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{returnErr: fmt.Errorf("job: %w", repository.ErrNotFound)}, nil).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPost, "/api/research/jobs/nonexistent/cancel", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestCancelResearchJobAlreadyTerminal(t *testing.T) {
	job := &models.ResearchJob{
		ID:     "job-1",
		Status: models.ResearchStatusPublished,
	}

	r := chi.NewRouter()
	handler.NewResearchHandler(&mockResearchRepo{job: job}, nil).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPost, "/api/research/jobs/job-1/cancel", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
