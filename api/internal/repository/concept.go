package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sean/apollo/api/internal/models"
)

// ConceptRepository defines read operations for concepts.
type ConceptRepository interface {
	ListConcepts(ctx context.Context, params models.PaginationParams, topicID string) (*models.PaginatedResponse[models.ConceptSummary], error)
	GetConceptByID(ctx context.Context, id string) (*models.ConceptDetail, error)
	GetConceptReferences(ctx context.Context, id string) ([]models.ConceptReference, error)
}

// SQLiteConceptRepository implements ConceptRepository using SQLite.
type SQLiteConceptRepository struct {
	db *sql.DB
}

// NewConceptRepository creates a new SQLiteConceptRepository.
func NewConceptRepository(db *sql.DB) *SQLiteConceptRepository {
	return &SQLiteConceptRepository{db: db}
}

const listConceptsSQL = `
SELECT id, name, definition, COALESCE(difficulty, ''), status,
       COALESCE(defined_in_topic, ''), aliases
FROM concepts
ORDER BY name
LIMIT ? OFFSET ?
`

const listConceptsByTopicSQL = `
SELECT id, name, definition, COALESCE(difficulty, ''), status,
       COALESCE(defined_in_topic, ''), aliases
FROM concepts
WHERE defined_in_topic = ?
ORDER BY name
LIMIT ? OFFSET ?
`

const countConceptsSQL = `SELECT COUNT(*) FROM concepts`

const countConceptsByTopicSQL = `SELECT COUNT(*) FROM concepts WHERE defined_in_topic = ?`

func (r *SQLiteConceptRepository) ListConcepts(ctx context.Context, params models.PaginationParams, topicID string) (*models.PaginatedResponse[models.ConceptSummary], error) {
	var total int

	var countErr error
	if topicID != "" {
		countErr = r.db.QueryRowContext(ctx, countConceptsByTopicSQL, topicID).Scan(&total)
	} else {
		countErr = r.db.QueryRowContext(ctx, countConceptsSQL).Scan(&total)
	}

	if countErr != nil {
		return nil, fmt.Errorf("count concepts: %w", countErr)
	}

	var rows *sql.Rows
	var queryErr error

	if topicID != "" {
		rows, queryErr = r.db.QueryContext(ctx, listConceptsByTopicSQL, topicID, params.PerPage, params.Offset())
	} else {
		rows, queryErr = r.db.QueryContext(ctx, listConceptsSQL, params.PerPage, params.Offset())
	}

	if queryErr != nil {
		return nil, fmt.Errorf("query concepts: %w", queryErr)
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

	if concepts == nil {
		concepts = []models.ConceptSummary{}
	}

	return &models.PaginatedResponse[models.ConceptSummary]{
		Items:   concepts,
		Total:   total,
		Page:    params.Page,
		PerPage: params.PerPage,
	}, nil
}

const getConceptSQL = `
SELECT id, name, definition, COALESCE(defined_in_lesson, ''), COALESCE(defined_in_topic, ''),
       COALESCE(difficulty, ''), COALESCE(flashcard_front, ''), COALESCE(flashcard_back, ''),
       status, aliases
FROM concepts
WHERE id = ?
`

const getRefsForConceptSQL = `
SELECT cr.lesson_id, l.title, COALESCE(cr.context, '')
FROM concept_references cr
JOIN lessons l ON l.id = cr.lesson_id
WHERE cr.concept_id = ?
`

func (r *SQLiteConceptRepository) GetConceptByID(ctx context.Context, id string) (*models.ConceptDetail, error) {
	cd := &models.ConceptDetail{}
	var aliasesRaw *string

	err := r.db.QueryRowContext(ctx, getConceptSQL, id).Scan(
		&cd.ID, &cd.Name, &cd.Definition, &cd.DefinedInLesson, &cd.DefinedInTopic,
		&cd.Difficulty, &cd.FlashcardFront, &cd.FlashcardBack,
		&cd.Status, &aliasesRaw,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query concept %s: %w", id, err)
	}

	cd.Aliases = models.ParseJSONStringSlice(aliasesRaw)

	refs, err := r.queryReferences(ctx, id)
	if err != nil {
		return nil, err
	}

	cd.References = refs

	return cd, nil
}

func (r *SQLiteConceptRepository) queryReferences(ctx context.Context, conceptID string) ([]models.ConceptReference, error) {
	rows, err := r.db.QueryContext(ctx, getRefsForConceptSQL, conceptID)
	if err != nil {
		return nil, fmt.Errorf("query references for concept %s: %w", conceptID, err)
	}
	defer rows.Close()

	var refs []models.ConceptReference

	for rows.Next() {
		var cr models.ConceptReference
		if err := rows.Scan(&cr.LessonID, &cr.LessonTitle, &cr.Context); err != nil {
			return nil, fmt.Errorf("scan concept reference: %w", err)
		}

		refs = append(refs, cr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate references: %w", err)
	}

	if refs == nil {
		refs = []models.ConceptReference{}
	}

	return refs, nil
}

func (r *SQLiteConceptRepository) GetConceptReferences(ctx context.Context, id string) ([]models.ConceptReference, error) {
	// First check concept exists.
	var exists bool

	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM concepts WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("check concept exists %s: %w", id, err)
	}

	if !exists {
		return nil, nil
	}

	return r.queryReferences(ctx, id)
}

// Verify interface compliance at compile time.
var _ ConceptRepository = (*SQLiteConceptRepository)(nil)
