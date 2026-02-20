package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sean/apollo/api/internal/models"
)

// ResearchJobRepository defines operations for research job persistence.
type ResearchJobRepository interface {
	CreateJob(ctx context.Context, input models.CreateResearchJobInput) (*models.ResearchJob, error)
	GetJobByID(ctx context.Context, id string) (*models.ResearchJob, error)
	ListJobs(ctx context.Context, params models.PaginationParams) (*models.PaginatedResponse[models.ResearchJobSummary], error)
	FindOldestByStatus(ctx context.Context, status string) (string, error)
	UpdateJobStatus(ctx context.Context, id string, status string, errorMsg string) error
	UpdateJobProgress(ctx context.Context, id string, progress models.ResearchProgress) error
	UpdateJobCurrentTopic(ctx context.Context, id string, topic string) error
}

// SQLiteResearchJobRepository implements ResearchJobRepository using SQLite.
type SQLiteResearchJobRepository struct {
	db *sql.DB
}

// NewResearchJobRepository creates a new SQLiteResearchJobRepository.
func NewResearchJobRepository(db *sql.DB) *SQLiteResearchJobRepository {
	return &SQLiteResearchJobRepository{db: db}
}

const createJobSQL = `
INSERT INTO research_jobs (id, root_topic, current_topic, status)
VALUES (?, ?, ?, ?)
`

const getJobByIDSQL = `
SELECT id, COALESCE(root_topic, ''), COALESCE(current_topic, ''), status,
       COALESCE(progress, ''), COALESCE(error, ''),
       COALESCE(started_at, ''), COALESCE(completed_at, '')
FROM research_jobs
WHERE id = ?
`

func (r *SQLiteResearchJobRepository) CreateJob(ctx context.Context, input models.CreateResearchJobInput) (*models.ResearchJob, error) {
	id := uuid.New().String()

	_, err := r.db.ExecContext(ctx, createJobSQL,
		id, input.Topic, input.Topic, models.ResearchStatusQueued,
	)
	if err != nil {
		return nil, classifyError(err, "create research job")
	}

	return r.GetJobByID(ctx, id)
}

func (r *SQLiteResearchJobRepository) GetJobByID(ctx context.Context, id string) (*models.ResearchJob, error) {
	job := &models.ResearchJob{}
	var progressStr, errStr string

	err := r.db.QueryRowContext(ctx, getJobByIDSQL, id).Scan(
		&job.ID, &job.RootTopic, &job.CurrentTopic, &job.Status,
		&progressStr, &errStr,
		&job.StartedAt, &job.CompletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("research job %s: %w", id, ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("get research job %s: %w", id, err)
	}

	if progressStr != "" {
		job.Progress = json.RawMessage(progressStr)
	}

	job.Error = errStr

	return job, nil
}

const countJobsSQL = `SELECT COUNT(*) FROM research_jobs`

const listJobsSQL = `
SELECT id, COALESCE(root_topic, ''), status,
       COALESCE(started_at, ''), COALESCE(completed_at, '')
FROM research_jobs
ORDER BY rowid DESC
LIMIT ? OFFSET ?
`

func (r *SQLiteResearchJobRepository) ListJobs(ctx context.Context, params models.PaginationParams) (*models.PaginatedResponse[models.ResearchJobSummary], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, countJobsSQL).Scan(&total); err != nil {
		return nil, fmt.Errorf("count research jobs: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, listJobsSQL, params.PerPage, params.Offset())
	if err != nil {
		return nil, fmt.Errorf("list research jobs: %w", err)
	}
	defer rows.Close()

	var items []models.ResearchJobSummary

	for rows.Next() {
		var job models.ResearchJobSummary
		if err := rows.Scan(&job.ID, &job.RootTopic, &job.Status, &job.StartedAt, &job.CompletedAt); err != nil {
			return nil, fmt.Errorf("scan research job summary: %w", err)
		}

		items = append(items, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate research jobs: %w", err)
	}

	if items == nil {
		items = []models.ResearchJobSummary{}
	}

	return &models.PaginatedResponse[models.ResearchJobSummary]{
		Items:   items,
		Total:   total,
		Page:    params.Page,
		PerPage: params.PerPage,
	}, nil
}

const updateJobStatusSQL = `
UPDATE research_jobs
SET status = ?,
    error = ?,
    started_at = CASE
      WHEN started_at IS NULL AND ? != 'queued' THEN ?
      ELSE started_at
    END,
    completed_at = CASE
      WHEN ? IN ('published', 'failed', 'cancelled') THEN ?
      ELSE completed_at
    END
WHERE id = ?
`

func (r *SQLiteResearchJobRepository) UpdateJobStatus(ctx context.Context, id string, status string, errorMsg string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	result, err := r.db.ExecContext(ctx, updateJobStatusSQL,
		status, nullIfEmpty(errorMsg),
		status, now,
		status, now,
		id,
	)
	if err != nil {
		return classifyError(err, "update research job status")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update research job status rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("research job %s: %w", id, ErrNotFound)
	}

	return nil
}

const updateJobProgressSQL = `UPDATE research_jobs SET progress = ? WHERE id = ?`

func (r *SQLiteResearchJobRepository) UpdateJobProgress(ctx context.Context, id string, progress models.ResearchProgress) error {
	progressJSON, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("marshal research progress: %w", err)
	}

	result, err := r.db.ExecContext(ctx, updateJobProgressSQL, string(progressJSON), id)
	if err != nil {
		return fmt.Errorf("update research job progress: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update research job progress rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("research job %s: %w", id, ErrNotFound)
	}

	return nil
}

const updateJobCurrentTopicSQL = `UPDATE research_jobs SET current_topic = ? WHERE id = ?`

func (r *SQLiteResearchJobRepository) UpdateJobCurrentTopic(ctx context.Context, id string, topic string) error {
	result, err := r.db.ExecContext(ctx, updateJobCurrentTopicSQL, topic, id)
	if err != nil {
		return fmt.Errorf("update research job current topic: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update research job current topic rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("research job %s: %w", id, ErrNotFound)
	}

	return nil
}

const findOldestByStatusSQL = `
SELECT id FROM research_jobs WHERE status = ? ORDER BY rowid ASC LIMIT 1
`

// FindOldestByStatus returns the ID of the oldest job with the given status,
// or "" if none found.
func (r *SQLiteResearchJobRepository) FindOldestByStatus(ctx context.Context, status string) (string, error) {
	var id string

	err := r.db.QueryRowContext(ctx, findOldestByStatusSQL, status).Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("find oldest %s job: %w", status, err)
	}

	return id, nil
}

// Verify interface compliance at compile time.
var _ ResearchJobRepository = (*SQLiteResearchJobRepository)(nil)
