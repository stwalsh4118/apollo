package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sean/apollo/api/internal/models"
)

// TopicRepository defines read operations for topics.
type TopicRepository interface {
	ListTopics(ctx context.Context) ([]models.TopicSummary, error)
	GetTopicByID(ctx context.Context, id string) (*models.TopicDetail, error)
	GetTopicFull(ctx context.Context, id string) (*models.TopicFull, error)
}

// SQLiteTopicRepository implements TopicRepository using SQLite.
type SQLiteTopicRepository struct {
	db *sql.DB
}

// NewTopicRepository creates a new SQLiteTopicRepository.
func NewTopicRepository(db *sql.DB) *SQLiteTopicRepository {
	return &SQLiteTopicRepository{db: db}
}

const listTopicsSQL = `
SELECT t.id, t.title, COALESCE(t.description, ''), COALESCE(t.difficulty, ''),
       COALESCE(t.estimated_hours, 0), t.tags, t.status,
       (SELECT COUNT(*) FROM modules m WHERE m.topic_id = t.id) AS module_count
FROM topics t
ORDER BY t.title
`

func (r *SQLiteTopicRepository) ListTopics(ctx context.Context) ([]models.TopicSummary, error) {
	rows, err := r.db.QueryContext(ctx, listTopicsSQL)
	if err != nil {
		return nil, fmt.Errorf("query topics: %w", err)
	}
	defer rows.Close()

	var topics []models.TopicSummary

	for rows.Next() {
		var ts models.TopicSummary
		var tagsRaw *string

		if err := rows.Scan(
			&ts.ID, &ts.Title, &ts.Description, &ts.Difficulty,
			&ts.EstimatedHours, &tagsRaw, &ts.Status, &ts.ModuleCount,
		); err != nil {
			return nil, fmt.Errorf("scan topic: %w", err)
		}

		ts.Tags = models.ParseJSONStringSlice(tagsRaw)
		topics = append(topics, ts)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate topics: %w", err)
	}

	if topics == nil {
		topics = []models.TopicSummary{}
	}

	return topics, nil
}

const getTopicSQL = `
SELECT id, title, COALESCE(description, ''), COALESCE(difficulty, ''),
       COALESCE(estimated_hours, 0), tags, status, version,
       source_urls, COALESCE(generated_at, ''), COALESCE(generated_by, ''),
       COALESCE(parent_topic_id, ''), created_at, updated_at
FROM topics
WHERE id = ?
`

const getModulesForTopicSQL = `
SELECT id, title, COALESCE(description, ''), COALESCE(estimated_minutes, 0), sort_order
FROM modules
WHERE topic_id = ?
ORDER BY sort_order
`

func (r *SQLiteTopicRepository) GetTopicByID(ctx context.Context, id string) (*models.TopicDetail, error) {
	td := &models.TopicDetail{}
	var tagsRaw, sourceURLsRaw *string

	err := r.db.QueryRowContext(ctx, getTopicSQL, id).Scan(
		&td.ID, &td.Title, &td.Description, &td.Difficulty,
		&td.EstimatedHours, &tagsRaw, &td.Status, &td.Version,
		&sourceURLsRaw, &td.GeneratedAt, &td.GeneratedBy,
		&td.ParentTopicID, &td.CreatedAt, &td.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query topic %s: %w", id, err)
	}

	td.Tags = models.ParseJSONStringSlice(tagsRaw)
	td.SourceURLs = models.ParseJSONStringSlice(sourceURLsRaw)

	modules, err := r.queryModuleSummaries(ctx, id)
	if err != nil {
		return nil, err
	}

	td.Modules = modules

	return td, nil
}

func (r *SQLiteTopicRepository) queryModuleSummaries(ctx context.Context, topicID string) ([]models.ModuleSummary, error) {
	rows, err := r.db.QueryContext(ctx, getModulesForTopicSQL, topicID)
	if err != nil {
		return nil, fmt.Errorf("query modules for topic %s: %w", topicID, err)
	}
	defer rows.Close()

	var modules []models.ModuleSummary

	for rows.Next() {
		var ms models.ModuleSummary
		if err := rows.Scan(&ms.ID, &ms.Title, &ms.Description, &ms.EstimatedMinutes, &ms.SortOrder); err != nil {
			return nil, fmt.Errorf("scan module: %w", err)
		}

		modules = append(modules, ms)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate modules: %w", err)
	}

	if modules == nil {
		modules = []models.ModuleSummary{}
	}

	return modules, nil
}

const getModulesFullSQL = `
SELECT id, topic_id, title, COALESCE(description, ''), learning_objectives,
       COALESCE(estimated_minutes, 0), sort_order, assessment
FROM modules
WHERE topic_id = ?
ORDER BY sort_order
`

const getLessonsForModuleSQL = `
SELECT id, module_id, title, sort_order, COALESCE(estimated_minutes, 0),
       content, examples, exercises, review_questions
FROM lessons
WHERE module_id = ?
ORDER BY sort_order
`

const getConceptsForLessonSQL = `
SELECT c.id, c.name, c.definition, COALESCE(c.difficulty, ''), c.status,
       COALESCE(c.defined_in_topic, ''), c.aliases
FROM concept_references cr
JOIN concepts c ON c.id = cr.concept_id
WHERE cr.lesson_id = ?
`

func (r *SQLiteTopicRepository) GetTopicFull(ctx context.Context, id string) (*models.TopicFull, error) {
	tf := &models.TopicFull{}
	var tagsRaw, sourceURLsRaw *string

	err := r.db.QueryRowContext(ctx, getTopicSQL, id).Scan(
		&tf.ID, &tf.Title, &tf.Description, &tf.Difficulty,
		&tf.EstimatedHours, &tagsRaw, &tf.Status, &tf.Version,
		&sourceURLsRaw, &tf.GeneratedAt, &tf.GeneratedBy,
		&tf.ParentTopicID, &tf.CreatedAt, &tf.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query topic full %s: %w", id, err)
	}

	tf.Tags = models.ParseJSONStringSlice(tagsRaw)
	tf.SourceURLs = models.ParseJSONStringSlice(sourceURLsRaw)

	modules, err := r.queryModulesFull(ctx, id)
	if err != nil {
		return nil, err
	}

	tf.Modules = modules

	return tf, nil
}

func (r *SQLiteTopicRepository) queryModulesFull(ctx context.Context, topicID string) ([]models.ModuleFull, error) {
	rows, err := r.db.QueryContext(ctx, getModulesFullSQL, topicID)
	if err != nil {
		return nil, fmt.Errorf("query modules full for topic %s: %w", topicID, err)
	}

	var modules []models.ModuleFull

	for rows.Next() {
		var mf models.ModuleFull
		var loRaw, assessRaw *string

		if err := rows.Scan(
			&mf.ID, &mf.TopicID, &mf.Title, &mf.Description, &loRaw,
			&mf.EstimatedMinutes, &mf.SortOrder, &assessRaw,
		); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan module full: %w", err)
		}

		mf.LearningObjectives = models.ParseJSONStringSlice(loRaw)
		mf.Assessment = models.ParseJSONRaw(assessRaw)
		modules = append(modules, mf)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate modules full: %w", err)
	}

	rows.Close()

	// Connection is now released — safe to make nested queries.
	for i := range modules {
		lessons, err := r.queryLessonsFull(ctx, modules[i].ID)
		if err != nil {
			return nil, err
		}

		modules[i].Lessons = lessons
	}

	if modules == nil {
		modules = []models.ModuleFull{}
	}

	return modules, nil
}

func (r *SQLiteTopicRepository) queryLessonsFull(ctx context.Context, moduleID string) ([]models.LessonFull, error) {
	rows, err := r.db.QueryContext(ctx, getLessonsForModuleSQL, moduleID)
	if err != nil {
		return nil, fmt.Errorf("query lessons for module %s: %w", moduleID, err)
	}

	var lessons []models.LessonFull

	for rows.Next() {
		var lf models.LessonFull
		var contentRaw, examplesRaw, exercisesRaw, reviewRaw *string

		if err := rows.Scan(
			&lf.ID, &lf.ModuleID, &lf.Title, &lf.SortOrder, &lf.EstimatedMinutes,
			&contentRaw, &examplesRaw, &exercisesRaw, &reviewRaw,
		); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan lesson full: %w", err)
		}

		lf.Content = models.ParseJSONRaw(contentRaw)
		lf.Examples = models.ParseJSONRaw(examplesRaw)
		lf.Exercises = models.ParseJSONRaw(exercisesRaw)
		lf.ReviewQuestions = models.ParseJSONRaw(reviewRaw)
		lessons = append(lessons, lf)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate lessons full: %w", err)
	}

	rows.Close()

	// Connection is now released — safe to make nested queries.
	for i := range lessons {
		concepts, err := r.queryConceptsForLesson(ctx, lessons[i].ID)
		if err != nil {
			return nil, err
		}

		lessons[i].Concepts = concepts
	}

	if lessons == nil {
		lessons = []models.LessonFull{}
	}

	return lessons, nil
}

func (r *SQLiteTopicRepository) queryConceptsForLesson(ctx context.Context, lessonID string) ([]models.ConceptSummary, error) {
	rows, err := r.db.QueryContext(ctx, getConceptsForLessonSQL, lessonID)
	if err != nil {
		return nil, fmt.Errorf("query concepts for lesson %s: %w", lessonID, err)
	}
	defer rows.Close()

	var concepts []models.ConceptSummary

	for rows.Next() {
		var cs models.ConceptSummary
		var aliasesRaw *string

		if err := rows.Scan(&cs.ID, &cs.Name, &cs.Definition, &cs.Difficulty, &cs.Status, &cs.DefinedInTopic, &aliasesRaw); err != nil {
			return nil, fmt.Errorf("scan concept: %w", err)
		}

		cs.Aliases = models.ParseJSONStringSlice(aliasesRaw)
		concepts = append(concepts, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate concepts: %w", err)
	}

	return concepts, nil
}

// Verify interface compliance at compile time.
var _ TopicRepository = (*SQLiteTopicRepository)(nil)
