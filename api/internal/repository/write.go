package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sean/apollo/api/internal/models"
)

// WriteRepository defines write operations for curriculum data.
type WriteRepository interface {
	CreateTopic(ctx context.Context, input models.TopicInput) error
	UpdateTopic(ctx context.Context, id string, input models.TopicInput) error
	CreateModule(ctx context.Context, input models.ModuleInput) error
	CreateLesson(ctx context.Context, input models.LessonInput) error
	CreateConcept(ctx context.Context, input models.ConceptInput) error
	CreateConceptReference(ctx context.Context, conceptID string, input models.ConceptReferenceInput) error
	CreatePrerequisite(ctx context.Context, input models.PrerequisiteInput) error
	CreateRelation(ctx context.Context, input models.RelationInput) error
}

// SQLiteWriteRepository implements WriteRepository using SQLite.
type SQLiteWriteRepository struct {
	db *sql.DB
}

// NewWriteRepository creates a new SQLiteWriteRepository.
func NewWriteRepository(db *sql.DB) *SQLiteWriteRepository {
	return &SQLiteWriteRepository{db: db}
}

const createTopicSQL = `
INSERT INTO topics (id, title, description, difficulty, estimated_hours, tags, status, version,
                    source_urls, generated_at, generated_by, parent_topic_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (r *SQLiteWriteRepository) CreateTopic(ctx context.Context, input models.TopicInput) error {
	tags := marshalJSONOrNil(input.Tags)
	sourceURLs := marshalJSONOrNil(input.SourceURLs)
	version := input.Version
	if version == 0 {
		version = 1
	}

	_, err := r.db.ExecContext(ctx, createTopicSQL,
		input.ID, input.Title, nullIfEmpty(input.Description), nullIfEmpty(input.Difficulty),
		nullIfZeroFloat(input.EstimatedHours), tags, input.Status, version,
		sourceURLs, nullIfEmpty(input.GeneratedAt), nullIfEmpty(input.GeneratedBy),
		nullIfEmpty(input.ParentTopicID),
	)
	if err != nil {
		return classifyError(err, "create topic")
	}

	return r.upsertSearchIndex(ctx, "topic", input.ID, input.Title, input.Description)
}

const updateTopicSQL = `
UPDATE topics SET title = ?, description = ?, difficulty = ?, estimated_hours = ?,
                  tags = ?, status = ?, version = ?,
                  source_urls = ?, generated_at = ?, generated_by = ?,
                  parent_topic_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
`

func (r *SQLiteWriteRepository) UpdateTopic(ctx context.Context, id string, input models.TopicInput) error {
	tags := marshalJSONOrNil(input.Tags)
	sourceURLs := marshalJSONOrNil(input.SourceURLs)
	version := input.Version
	if version == 0 {
		version = 1
	}

	result, err := r.db.ExecContext(ctx, updateTopicSQL,
		input.Title, nullIfEmpty(input.Description), nullIfEmpty(input.Difficulty),
		nullIfZeroFloat(input.EstimatedHours), tags, input.Status, version,
		sourceURLs, nullIfEmpty(input.GeneratedAt), nullIfEmpty(input.GeneratedBy),
		nullIfEmpty(input.ParentTopicID), id,
	)
	if err != nil {
		return classifyError(err, "update topic")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update topic rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("topic %s: %w", id, ErrNotFound)
	}

	return r.upsertSearchIndex(ctx, "topic", id, input.Title, input.Description)
}

const createModuleSQL = `
INSERT INTO modules (id, topic_id, title, description, learning_objectives, estimated_minutes, sort_order, assessment)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

func (r *SQLiteWriteRepository) CreateModule(ctx context.Context, input models.ModuleInput) error {
	lo := marshalJSONOrNil(input.LearningObjectives)
	assessment := rawJSONOrNil(input.Assessment)

	_, err := r.db.ExecContext(ctx, createModuleSQL,
		input.ID, input.TopicID, input.Title, nullIfEmpty(input.Description),
		lo, nullIfZero(input.EstimatedMinutes), input.SortOrder, assessment,
	)
	if err != nil {
		return classifyError(err, "create module")
	}

	return nil
}

const createLessonSQL = `
INSERT INTO lessons (id, module_id, title, sort_order, estimated_minutes, content, examples, exercises, review_questions)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (r *SQLiteWriteRepository) CreateLesson(ctx context.Context, input models.LessonInput) error {
	_, err := r.db.ExecContext(ctx, createLessonSQL,
		input.ID, input.ModuleID, input.Title, input.SortOrder,
		nullIfZero(input.EstimatedMinutes), string(input.Content),
		rawJSONOrNil(input.Examples), rawJSONOrNil(input.Exercises),
		rawJSONOrNil(input.ReviewQuestions),
	)
	if err != nil {
		return classifyError(err, "create lesson")
	}

	return r.upsertSearchIndex(ctx, "lesson", input.ID, input.Title, "")
}

const createConceptSQL = `
INSERT INTO concepts (id, name, definition, defined_in_lesson, defined_in_topic,
                      difficulty, flashcard_front, flashcard_back, status, aliases)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (r *SQLiteWriteRepository) CreateConcept(ctx context.Context, input models.ConceptInput) error {
	aliases := marshalJSONOrNil(input.Aliases)
	status := input.Status
	if status == "" {
		status = "active"
	}

	_, err := r.db.ExecContext(ctx, createConceptSQL,
		input.ID, input.Name, input.Definition,
		nullIfEmpty(input.DefinedInLesson), nullIfEmpty(input.DefinedInTopic),
		nullIfEmpty(input.Difficulty), nullIfEmpty(input.FlashcardFront),
		nullIfEmpty(input.FlashcardBack), status, aliases,
	)
	if err != nil {
		return classifyError(err, "create concept")
	}

	return r.upsertSearchIndex(ctx, "concept", input.ID, input.Name, input.Definition)
}

const createConceptRefSQL = `
INSERT INTO concept_references (concept_id, lesson_id, context)
VALUES (?, ?, ?)
`

func (r *SQLiteWriteRepository) CreateConceptReference(ctx context.Context, conceptID string, input models.ConceptReferenceInput) error {
	_, err := r.db.ExecContext(ctx, createConceptRefSQL,
		conceptID, input.LessonID, nullIfEmpty(input.Context),
	)
	if err != nil {
		return classifyError(err, "create concept reference")
	}

	return nil
}

const createPrerequisiteSQL = `
INSERT INTO topic_prerequisites (topic_id, prerequisite_topic_id, priority, reason)
VALUES (?, ?, ?, ?)
`

func (r *SQLiteWriteRepository) CreatePrerequisite(ctx context.Context, input models.PrerequisiteInput) error {
	_, err := r.db.ExecContext(ctx, createPrerequisiteSQL,
		input.TopicID, input.PrerequisiteTopicID, input.Priority, nullIfEmpty(input.Reason),
	)
	if err != nil {
		return classifyError(err, "create prerequisite")
	}

	return nil
}

const createRelationSQL = `
INSERT INTO topic_relations (topic_a, topic_b, relation_type, description)
VALUES (?, ?, ?, ?)
`

func (r *SQLiteWriteRepository) CreateRelation(ctx context.Context, input models.RelationInput) error {
	_, err := r.db.ExecContext(ctx, createRelationSQL,
		input.TopicA, input.TopicB, input.RelationType, nullIfEmpty(input.Description),
	)
	if err != nil {
		return classifyError(err, "create relation")
	}

	return nil
}

// FTS5 virtual tables don't support ON CONFLICT. Use delete + insert instead.
const deleteSearchSQL = `DELETE FROM search_index WHERE entity_type = ? AND entity_id = ?`
const insertSearchSQL = `INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)`

func (r *SQLiteWriteRepository) upsertSearchIndex(ctx context.Context, entityType, entityID, title, body string) error {
	if _, err := r.db.ExecContext(ctx, deleteSearchSQL, entityType, entityID); err != nil {
		return fmt.Errorf("delete search index for %s %s: %w", entityType, entityID, err)
	}

	_, err := r.db.ExecContext(ctx, insertSearchSQL, entityType, entityID, title, body)
	if err != nil {
		return fmt.Errorf("insert search index for %s %s: %w", entityType, entityID, err)
	}

	return nil
}

func classifyError(err error, op string) error {
	msg := err.Error()

	if strings.Contains(msg, "UNIQUE constraint failed") {
		return fmt.Errorf("%s: %w", op, ErrDuplicate)
	}

	if strings.Contains(msg, "FOREIGN KEY constraint failed") {
		return fmt.Errorf("%s: %w", op, ErrFKViolation)
	}

	if strings.Contains(msg, "CHECK constraint failed") {
		return fmt.Errorf("%s: %w", op, ErrCheckViolation)
	}

	return fmt.Errorf("%s: %w", op, err)
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}

	return s
}

func nullIfZero(n int) any {
	if n == 0 {
		return nil
	}

	return n
}

func nullIfZeroFloat(f float64) any {
	if f == 0 {
		return nil
	}

	return f
}

func marshalJSONOrNil(slice []string) any {
	if len(slice) == 0 {
		return nil
	}

	b, err := json.Marshal(slice)
	if err != nil {
		return nil
	}

	return string(b)
}

func rawJSONOrNil(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}

	return string(raw)
}

// Verify interface compliance at compile time.
var _ WriteRepository = (*SQLiteWriteRepository)(nil)
