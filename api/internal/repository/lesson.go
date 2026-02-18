package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sean/apollo/api/internal/models"
)

// LessonRepository defines read operations for lessons.
type LessonRepository interface {
	GetLessonByID(ctx context.Context, id string) (*models.LessonDetail, error)
}

// SQLiteLessonRepository implements LessonRepository using SQLite.
type SQLiteLessonRepository struct {
	db *sql.DB
}

// NewLessonRepository creates a new SQLiteLessonRepository.
func NewLessonRepository(db *sql.DB) *SQLiteLessonRepository {
	return &SQLiteLessonRepository{db: db}
}

const getLessonSQL = `
SELECT id, module_id, title, sort_order, COALESCE(estimated_minutes, 0),
       content, examples, exercises, review_questions
FROM lessons
WHERE id = ?
`

func (r *SQLiteLessonRepository) GetLessonByID(ctx context.Context, id string) (*models.LessonDetail, error) {
	ld := &models.LessonDetail{}
	var contentRaw, examplesRaw, exercisesRaw, reviewRaw *string

	err := r.db.QueryRowContext(ctx, getLessonSQL, id).Scan(
		&ld.ID, &ld.ModuleID, &ld.Title, &ld.SortOrder, &ld.EstimatedMinutes,
		&contentRaw, &examplesRaw, &exercisesRaw, &reviewRaw,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query lesson %s: %w", id, err)
	}

	ld.Content = models.ParseJSONRaw(contentRaw)
	ld.Examples = models.ParseJSONRaw(examplesRaw)
	ld.Exercises = models.ParseJSONRaw(exercisesRaw)
	ld.ReviewQuestions = models.ParseJSONRaw(reviewRaw)

	return ld, nil
}

var _ LessonRepository = (*SQLiteLessonRepository)(nil)
