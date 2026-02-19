package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sean/apollo/api/internal/models"
)

// ProgressRepository defines operations for learning progress tracking.
type ProgressRepository interface {
	GetTopicProgress(ctx context.Context, topicID string) (*models.TopicProgress, error)
	UpdateLessonProgress(ctx context.Context, lessonID string, input models.UpdateProgressInput) (*models.LessonProgress, error)
	GetProgressSummary(ctx context.Context) (*models.ProgressSummary, error)
}

// SQLiteProgressRepository implements ProgressRepository using SQLite.
type SQLiteProgressRepository struct {
	db *sql.DB
}

// NewProgressRepository creates a new SQLiteProgressRepository.
func NewProgressRepository(db *sql.DB) *SQLiteProgressRepository {
	return &SQLiteProgressRepository{db: db}
}

const checkTopicExistsSQL = `SELECT EXISTS(SELECT 1 FROM topics WHERE id = ?)`

const getTopicProgressSQL = `
SELECT l.id, l.title, COALESCE(lp.status, 'not_started'),
       COALESCE(lp.started_at, ''), COALESCE(lp.completed_at, ''),
       COALESCE(lp.notes, '')
FROM modules m
JOIN lessons l ON l.module_id = m.id
LEFT JOIN learning_progress lp ON lp.lesson_id = l.id
WHERE m.topic_id = ?
ORDER BY m.sort_order, l.sort_order
`

func (r *SQLiteProgressRepository) GetTopicProgress(ctx context.Context, topicID string) (*models.TopicProgress, error) {
	// Verify topic exists.
	var exists bool
	if err := r.db.QueryRowContext(ctx, checkTopicExistsSQL, topicID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check topic %s: %w", topicID, err)
	}

	if !exists {
		return nil, ErrNotFound
	}

	rows, err := r.db.QueryContext(ctx, getTopicProgressSQL, topicID)
	if err != nil {
		return nil, fmt.Errorf("query topic progress %s: %w", topicID, err)
	}
	defer rows.Close()

	tp := &models.TopicProgress{TopicID: topicID}

	for rows.Next() {
		var lp models.LessonProgress
		if err := rows.Scan(&lp.LessonID, &lp.LessonTitle, &lp.Status, &lp.StartedAt, &lp.CompletedAt, &lp.Notes); err != nil {
			return nil, fmt.Errorf("scan lesson progress: %w", err)
		}

		tp.Lessons = append(tp.Lessons, lp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate topic progress: %w", err)
	}

	if tp.Lessons == nil {
		tp.Lessons = []models.LessonProgress{}
	}

	return tp, nil
}

const checkLessonExistsSQL = `SELECT EXISTS(SELECT 1 FROM lessons WHERE id = ?)`

const upsertProgressSQL = `
INSERT INTO learning_progress (lesson_id, status, started_at, completed_at, notes)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT(lesson_id) DO UPDATE SET
  status = excluded.status,
  started_at = COALESCE(learning_progress.started_at, excluded.started_at),
  completed_at = excluded.completed_at,
  notes = excluded.notes
`

const getProgressSQL = `
SELECT lesson_id, status, COALESCE(started_at, ''), COALESCE(completed_at, ''), COALESCE(notes, '')
FROM learning_progress
WHERE lesson_id = ?
`

func (r *SQLiteProgressRepository) UpdateLessonProgress(ctx context.Context, lessonID string, input models.UpdateProgressInput) (*models.LessonProgress, error) {
	// Verify lesson exists.
	var exists bool
	if err := r.db.QueryRowContext(ctx, checkLessonExistsSQL, lessonID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check lesson %s: %w", lessonID, err)
	}

	if !exists {
		return nil, ErrNotFound
	}

	now := time.Now().UTC().Format(time.RFC3339)

	var startedAt, completedAt string
	if input.Status != models.ProgressStatusNotStarted {
		startedAt = now
	}

	if input.Status == models.ProgressStatusCompleted {
		completedAt = now
	}

	if _, err := r.db.ExecContext(ctx, upsertProgressSQL, lessonID, input.Status, startedAt, completedAt, input.Notes); err != nil {
		return nil, fmt.Errorf("upsert progress %s: %w", lessonID, err)
	}

	// Read back the stored row to return accurate timestamps.
	lp := &models.LessonProgress{}
	if err := r.db.QueryRowContext(ctx, getProgressSQL, lessonID).Scan(
		&lp.LessonID, &lp.Status, &lp.StartedAt, &lp.CompletedAt, &lp.Notes,
	); err != nil {
		return nil, fmt.Errorf("read back progress %s: %w", lessonID, err)
	}

	return lp, nil
}

const getProgressSummarySQL = `
SELECT
  (SELECT COUNT(*) FROM lessons) AS total_lessons,
  (SELECT COUNT(*) FROM learning_progress WHERE status = 'completed') AS completed_lessons,
  (SELECT COUNT(DISTINCT m.topic_id)
   FROM learning_progress lp
   JOIN lessons l ON l.id = lp.lesson_id
   JOIN modules m ON m.id = l.module_id
   WHERE lp.status IN ('in_progress', 'completed')) AS active_topics
`

func (r *SQLiteProgressRepository) GetProgressSummary(ctx context.Context) (*models.ProgressSummary, error) {
	ps := &models.ProgressSummary{}

	if err := r.db.QueryRowContext(ctx, getProgressSummarySQL).Scan(
		&ps.TotalLessons, &ps.CompletedLessons, &ps.ActiveTopics,
	); err != nil {
		return nil, fmt.Errorf("query progress summary: %w", err)
	}

	if ps.TotalLessons > 0 {
		ps.CompletionPercentage = float64(ps.CompletedLessons) / float64(ps.TotalLessons) * 100
	}

	return ps, nil
}

// Verify interface compliance at compile time.
var _ ProgressRepository = (*SQLiteProgressRepository)(nil)
