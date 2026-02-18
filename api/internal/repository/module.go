package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sean/apollo/api/internal/models"
)

// ModuleRepository defines read operations for modules.
type ModuleRepository interface {
	GetModuleByID(ctx context.Context, id string) (*models.ModuleDetail, error)
}

// SQLiteModuleRepository implements ModuleRepository using SQLite.
type SQLiteModuleRepository struct {
	db *sql.DB
}

// NewModuleRepository creates a new SQLiteModuleRepository.
func NewModuleRepository(db *sql.DB) *SQLiteModuleRepository {
	return &SQLiteModuleRepository{db: db}
}

const getModuleSQL = `
SELECT id, topic_id, title, COALESCE(description, ''), learning_objectives,
       COALESCE(estimated_minutes, 0), sort_order, assessment
FROM modules
WHERE id = ?
`

const getLessonsForModuleSummarySQL = `
SELECT id, title, sort_order, COALESCE(estimated_minutes, 0)
FROM lessons
WHERE module_id = ?
ORDER BY sort_order
`

func (r *SQLiteModuleRepository) GetModuleByID(ctx context.Context, id string) (*models.ModuleDetail, error) {
	md := &models.ModuleDetail{}
	var loRaw, assessRaw *string

	err := r.db.QueryRowContext(ctx, getModuleSQL, id).Scan(
		&md.ID, &md.TopicID, &md.Title, &md.Description, &loRaw,
		&md.EstimatedMinutes, &md.SortOrder, &assessRaw,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query module %s: %w", id, err)
	}

	md.LearningObjectives = models.ParseJSONStringSlice(loRaw)
	md.Assessment = models.ParseJSONRaw(assessRaw)

	lessons, err := r.queryLessonSummaries(ctx, id)
	if err != nil {
		return nil, err
	}

	md.Lessons = lessons

	return md, nil
}

func (r *SQLiteModuleRepository) queryLessonSummaries(ctx context.Context, moduleID string) ([]models.LessonSummary, error) {
	rows, err := r.db.QueryContext(ctx, getLessonsForModuleSummarySQL, moduleID)
	if err != nil {
		return nil, fmt.Errorf("query lessons for module %s: %w", moduleID, err)
	}
	defer rows.Close()

	var lessons []models.LessonSummary

	for rows.Next() {
		var ls models.LessonSummary
		if err := rows.Scan(&ls.ID, &ls.Title, &ls.SortOrder, &ls.EstimatedMinutes); err != nil {
			return nil, fmt.Errorf("scan lesson summary: %w", err)
		}

		lessons = append(lessons, ls)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lessons: %w", err)
	}

	if lessons == nil {
		lessons = []models.LessonSummary{}
	}

	return lessons, nil
}

var _ ModuleRepository = (*SQLiteModuleRepository)(nil)
