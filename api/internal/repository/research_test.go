package repository_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
)

func TestCreateJob(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Go Concurrency",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	if job.ID == "" {
		t.Fatal("expected non-empty job ID")
	}

	if job.RootTopic != "Go Concurrency" {
		t.Fatalf("expected root_topic %q, got %q", "Go Concurrency", job.RootTopic)
	}

	if job.Status != models.ResearchStatusQueued {
		t.Fatalf("expected status %q, got %q", models.ResearchStatusQueued, job.Status)
	}

	if job.CurrentTopic != "Go Concurrency" {
		t.Fatalf("expected current_topic %q, got %q", "Go Concurrency", job.CurrentTopic)
	}
}

func TestGetJobByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	created, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Rust Lifetimes",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	retrieved, err := repo.GetJobByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Fatalf("expected ID %q, got %q", created.ID, retrieved.ID)
	}

	if retrieved.RootTopic != "Rust Lifetimes" {
		t.Fatalf("expected root_topic %q, got %q", "Rust Lifetimes", retrieved.RootTopic)
	}
}

func TestGetJobByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	_, err := repo.GetJobByID(context.Background(), "nonexistent")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListJobs(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	// Empty list.
	list, err := repo.ListJobs(context.Background(), models.PaginationParams{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}

	if list.Total != 0 {
		t.Fatalf("expected 0 total, got %d", list.Total)
	}

	if len(list.Items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(list.Items))
	}

	// Create some jobs.
	for _, topic := range []string{"Topic A", "Topic B", "Topic C"} {
		if _, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{Topic: topic}); err != nil {
			t.Fatalf("create job %q: %v", topic, err)
		}
	}

	list, err = repo.ListJobs(context.Background(), models.PaginationParams{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}

	if list.Total != 3 {
		t.Fatalf("expected 3 total, got %d", list.Total)
	}

	if len(list.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(list.Items))
	}

	// Most recent first (C, B, A).
	if list.Items[0].RootTopic != "Topic C" {
		t.Fatalf("expected first item to be 'Topic C', got %q", list.Items[0].RootTopic)
	}
}

func TestListJobsPagination(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	for _, topic := range []string{"A", "B", "C"} {
		if _, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{Topic: topic}); err != nil {
			t.Fatalf("create job: %v", err)
		}
	}

	list, err := repo.ListJobs(context.Background(), models.PaginationParams{Page: 1, PerPage: 2})
	if err != nil {
		t.Fatalf("list jobs page 1: %v", err)
	}

	if list.Total != 3 {
		t.Fatalf("expected 3 total, got %d", list.Total)
	}

	if len(list.Items) != 2 {
		t.Fatalf("expected 2 items on page 1, got %d", len(list.Items))
	}

	list, err = repo.ListJobs(context.Background(), models.PaginationParams{Page: 2, PerPage: 2})
	if err != nil {
		t.Fatalf("list jobs page 2: %v", err)
	}

	if len(list.Items) != 1 {
		t.Fatalf("expected 1 item on page 2, got %d", len(list.Items))
	}
}

func TestUpdateJobStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{Topic: "Test Topic"})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	// Transition to researching — sets started_at.
	if err := repo.UpdateJobStatus(context.Background(), job.ID, models.ResearchStatusResearching, ""); err != nil {
		t.Fatalf("update status to researching: %v", err)
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusResearching {
		t.Fatalf("expected status %q, got %q", models.ResearchStatusResearching, updated.Status)
	}

	if updated.StartedAt == "" {
		t.Fatal("expected started_at to be set")
	}

	// Transition to published — sets completed_at.
	if err := repo.UpdateJobStatus(context.Background(), job.ID, models.ResearchStatusPublished, ""); err != nil {
		t.Fatalf("update status to published: %v", err)
	}

	updated, err = repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusPublished {
		t.Fatalf("expected status %q, got %q", models.ResearchStatusPublished, updated.Status)
	}

	if updated.CompletedAt == "" {
		t.Fatal("expected completed_at to be set")
	}
}

func TestUpdateJobStatusFailed(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{Topic: "Fail Topic"})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	if err := repo.UpdateJobStatus(context.Background(), job.ID, models.ResearchStatusFailed, "pass 2 timed out"); err != nil {
		t.Fatalf("update status to failed: %v", err)
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusFailed {
		t.Fatalf("expected status %q, got %q", models.ResearchStatusFailed, updated.Status)
	}

	if updated.Error != "pass 2 timed out" {
		t.Fatalf("expected error %q, got %q", "pass 2 timed out", updated.Error)
	}

	if updated.CompletedAt == "" {
		t.Fatal("expected completed_at to be set on failure")
	}
}

func TestUpdateJobStatusNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	err := repo.UpdateJobStatus(context.Background(), "nonexistent", models.ResearchStatusFailed, "")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateJobProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{Topic: "Progress Topic"})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	progress := models.ResearchProgress{
		CurrentPass:      2,
		TotalPasses:      4,
		ModulesPlanned:   5,
		ModulesCompleted: 1,
		ConceptsFound:    8,
	}

	if err := repo.UpdateJobProgress(context.Background(), job.ID, progress); err != nil {
		t.Fatalf("update progress: %v", err)
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Progress == nil {
		t.Fatal("expected progress to be non-nil")
	}

	var decoded models.ResearchProgress
	if err := json.Unmarshal(updated.Progress, &decoded); err != nil {
		t.Fatalf("unmarshal progress: %v", err)
	}

	if decoded.CurrentPass != 2 {
		t.Fatalf("expected current_pass 2, got %d", decoded.CurrentPass)
	}

	if decoded.ConceptsFound != 8 {
		t.Fatalf("expected concepts_found 8, got %d", decoded.ConceptsFound)
	}
}

func TestUpdateJobProgressNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	err := repo.UpdateJobProgress(context.Background(), "nonexistent", models.ResearchProgress{})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateJobStatusStartedAtNotOverwritten(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{Topic: "Sticky Start"})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	// First transition sets started_at.
	if err := repo.UpdateJobStatus(context.Background(), job.ID, models.ResearchStatusResearching, ""); err != nil {
		t.Fatalf("update to researching: %v", err)
	}

	first, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	originalStartedAt := first.StartedAt
	if originalStartedAt == "" {
		t.Fatal("expected started_at to be set after first transition")
	}

	// Second transition should NOT overwrite started_at.
	if err := repo.UpdateJobStatus(context.Background(), job.ID, models.ResearchStatusResolving, ""); err != nil {
		t.Fatalf("update to resolving: %v", err)
	}

	second, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if second.StartedAt != originalStartedAt {
		t.Fatalf("expected started_at %q to be preserved, got %q", originalStartedAt, second.StartedAt)
	}
}

func TestUpdateJobCurrentTopic(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{Topic: "Original"})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	if err := repo.UpdateJobCurrentTopic(context.Background(), job.ID, "New Subtopic"); err != nil {
		t.Fatalf("update current topic: %v", err)
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.CurrentTopic != "New Subtopic" {
		t.Fatalf("expected current_topic %q, got %q", "New Subtopic", updated.CurrentTopic)
	}
}

func TestUpdateJobCurrentTopicNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)

	err := repo.UpdateJobCurrentTopic(context.Background(), "nonexistent", "topic")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
