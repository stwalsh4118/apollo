package research

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/sean/apollo/api/internal/schema"
)

// CurriculumIngester validates and stores a curriculum JSON output in SQLite.
type CurriculumIngester struct {
	db *sql.DB
}

// NewCurriculumIngester creates a new CurriculumIngester.
func NewCurriculumIngester(db *sql.DB) *CurriculumIngester {
	return &CurriculumIngester{db: db}
}

// Ingest validates rawJSON against the curriculum schema, parses it,
// and stores all entities in the database within a single transaction.
func (ing *CurriculumIngester) Ingest(ctx context.Context, rawJSON json.RawMessage) error {
	if err := schema.Validate(rawJSON); err != nil {
		return fmt.Errorf("schema validation: %w", err)
	}

	var curr CurriculumOutput
	if err := json.Unmarshal(rawJSON, &curr); err != nil {
		return fmt.Errorf("unmarshal curriculum: %w", err)
	}

	tx, err := ing.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	if err := ing.storeTopic(ctx, tx, &curr); err != nil {
		return err
	}

	for i, mod := range curr.Modules {
		if err := ing.storeModule(ctx, tx, curr.ID, &mod, i); err != nil {
			return err
		}

		for j, lesson := range mod.Lessons {
			if err := ing.storeLesson(ctx, tx, mod.ID, &lesson, j); err != nil {
				return err
			}

			for _, concept := range lesson.ConceptsTaught {
				if err := ing.storeConcept(ctx, tx, curr.ID, lesson.ID, &concept); err != nil {
					return err
				}
			}

			for _, ref := range lesson.ConceptsReferenced {
				if err := ing.storeConceptReference(ctx, tx, ref.ID, lesson.ID); err != nil {
					return err
				}
			}
		}
	}

	if err := ing.storePrerequisites(ctx, tx, curr.ID, &curr.Prerequisites); err != nil {
		return err
	}

	if err := ing.storeExpansionQueue(ctx, tx, curr.ID, &curr.Prerequisites); err != nil {
		return err
	}

	if err := ing.storeSearchIndex(ctx, tx, &curr); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

const insertTopicSQL = `
INSERT INTO topics (id, title, description, difficulty, estimated_hours, tags, status, version,
                    source_urls, generated_at, generated_by)
VALUES (?, ?, ?, ?, ?, ?, 'published', ?, ?, ?, 'research-agent')
`

func (ing *CurriculumIngester) storeTopic(ctx context.Context, tx *sql.Tx, curr *CurriculumOutput) error {
	tags, _ := json.Marshal(curr.Tags)
	sourceURLs, _ := json.Marshal(curr.SourceURLs)

	_, err := tx.ExecContext(ctx, insertTopicSQL,
		curr.ID, curr.Title, curr.Description, curr.Difficulty,
		curr.EstimatedHours, string(tags), curr.Version,
		string(sourceURLs), curr.GeneratedAt,
	)
	if err != nil {
		return fmt.Errorf("insert topic %s: %w", curr.ID, err)
	}

	return nil
}

const insertModuleSQL = `
INSERT INTO modules (id, topic_id, title, description, learning_objectives, estimated_minutes, sort_order, assessment)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

func (ing *CurriculumIngester) storeModule(ctx context.Context, tx *sql.Tx, topicID string, mod *ModuleOutput, idx int) error {
	lo, _ := json.Marshal(mod.LearningObjectives)
	sortOrder := mod.Order
	if sortOrder == 0 {
		sortOrder = idx + 1
	}

	var assessment any
	if len(mod.Assessment) > 0 {
		assessment = string(mod.Assessment)
	}

	_, err := tx.ExecContext(ctx, insertModuleSQL,
		mod.ID, topicID, mod.Title, mod.Description,
		string(lo), mod.EstimatedMinutes, sortOrder, assessment,
	)
	if err != nil {
		return fmt.Errorf("insert module %s: %w", mod.ID, err)
	}

	return nil
}

const insertLessonSQL = `
INSERT INTO lessons (id, module_id, title, sort_order, estimated_minutes, content, examples, exercises, review_questions)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (ing *CurriculumIngester) storeLesson(ctx context.Context, tx *sql.Tx, moduleID string, lesson *LessonOutput, idx int) error {
	sortOrder := lesson.Order
	if sortOrder == 0 {
		sortOrder = idx + 1
	}

	contentStr := rawOrNil(lesson.Content)
	examplesStr := rawOrNil(lesson.Examples)
	exercisesStr := rawOrNil(lesson.Exercises)
	reviewStr := rawOrNil(lesson.ReviewQuestions)

	_, err := tx.ExecContext(ctx, insertLessonSQL,
		lesson.ID, moduleID, lesson.Title, sortOrder,
		lesson.EstimatedMinutes, contentStr, examplesStr, exercisesStr, reviewStr,
	)
	if err != nil {
		return fmt.Errorf("insert lesson %s: %w", lesson.ID, err)
	}

	return nil
}

const insertConceptSQL = `
INSERT INTO concepts (id, name, definition, defined_in_lesson, defined_in_topic,
                      flashcard_front, flashcard_back, status)
VALUES (?, ?, ?, ?, ?, ?, ?, 'active')
`

func (ing *CurriculumIngester) storeConcept(ctx context.Context, tx *sql.Tx, topicID, lessonID string, concept *ConceptTaughtOut) error {
	_, err := tx.ExecContext(ctx, insertConceptSQL,
		concept.ID, concept.Name, concept.Definition,
		lessonID, topicID,
		concept.Flashcard.Front, concept.Flashcard.Back,
	)
	if err != nil {
		return fmt.Errorf("insert concept %s: %w", concept.ID, err)
	}

	// Also create a self-reference for the defining lesson.
	if err := ing.storeConceptReference(ctx, tx, concept.ID, lessonID); err != nil {
		return err
	}

	return nil
}

const insertConceptRefSQL = `
INSERT OR IGNORE INTO concept_references (concept_id, lesson_id, context)
VALUES (?, ?, '')
`

func (ing *CurriculumIngester) storeConceptReference(ctx context.Context, tx *sql.Tx, conceptID, lessonID string) error {
	_, err := tx.ExecContext(ctx, insertConceptRefSQL, conceptID, lessonID)
	if err != nil {
		return fmt.Errorf("insert concept reference %s -> %s: %w", conceptID, lessonID, err)
	}

	return nil
}

const checkTopicExistsSQL = `SELECT EXISTS(SELECT 1 FROM topics WHERE id = ?)`

const insertPrereqSQL = `
INSERT OR IGNORE INTO topic_prerequisites (topic_id, prerequisite_topic_id, priority, reason)
VALUES (?, ?, ?, ?)
`

func (ing *CurriculumIngester) storePrerequisites(ctx context.Context, tx *sql.Tx, topicID string, prereqs *PrerequisitesOutput) error {
	store := func(items []PrerequisiteItem, priority string) error {
		for _, item := range items {
			// Only store in topic_prerequisites if the prerequisite topic already exists.
			// Non-existent topics are handled via the expansion queue.
			var exists bool
			if err := tx.QueryRowContext(ctx, checkTopicExistsSQL, item.TopicID).Scan(&exists); err != nil {
				return fmt.Errorf("check prerequisite topic %s: %w", item.TopicID, err)
			}

			if !exists {
				continue
			}

			_, err := tx.ExecContext(ctx, insertPrereqSQL,
				topicID, item.TopicID, priority, item.Reason,
			)
			if err != nil {
				return fmt.Errorf("insert prerequisite %s -> %s: %w", topicID, item.TopicID, err)
			}
		}

		return nil
	}

	if err := store(prereqs.Essential, "essential"); err != nil {
		return err
	}

	if err := store(prereqs.Helpful, "helpful"); err != nil {
		return err
	}

	return store(prereqs.DeepBackground, "deep_background")
}

const insertExpansionQueueSQL = `
INSERT INTO expansion_queue (topic_id, requested_by_topic, priority, reason, status)
VALUES (?, ?, ?, ?, 'available')
`

func (ing *CurriculumIngester) storeExpansionQueue(ctx context.Context, tx *sql.Tx, topicID string, prereqs *PrerequisitesOutput) error {
	enqueue := func(items []PrerequisiteItem, priority string) error {
		for _, item := range items {
			_, err := tx.ExecContext(ctx, insertExpansionQueueSQL,
				item.TopicID, topicID, priority, item.Reason,
			)
			if err != nil {
				return fmt.Errorf("insert expansion queue %s: %w", item.TopicID, err)
			}
		}

		return nil
	}

	if err := enqueue(prereqs.Helpful, "helpful"); err != nil {
		return err
	}

	return enqueue(prereqs.DeepBackground, "deep_background")
}

const deleteSearchSQL = `DELETE FROM search_index WHERE entity_type = ? AND entity_id = ?`
const insertSearchSQL = `INSERT INTO search_index (entity_type, entity_id, title, body) VALUES (?, ?, ?, ?)`

func (ing *CurriculumIngester) storeSearchIndex(ctx context.Context, tx *sql.Tx, curr *CurriculumOutput) error {
	upsert := func(entityType, entityID, title, body string) error {
		if _, err := tx.ExecContext(ctx, deleteSearchSQL, entityType, entityID); err != nil {
			return fmt.Errorf("delete search index %s %s: %w", entityType, entityID, err)
		}

		_, err := tx.ExecContext(ctx, insertSearchSQL, entityType, entityID, title, body)
		if err != nil {
			return fmt.Errorf("insert search index %s %s: %w", entityType, entityID, err)
		}

		return nil
	}

	if err := upsert("topic", curr.ID, curr.Title, curr.Description); err != nil {
		return err
	}

	for _, mod := range curr.Modules {
		for _, lesson := range mod.Lessons {
			if err := upsert("lesson", lesson.ID, lesson.Title, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

func rawOrNil(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}

	return string(raw)
}

// IngestResult holds counts from a successful ingestion.
type IngestResult struct {
	ModulesCreated  int
	LessonsCreated  int
	ConceptsCreated int
}

// IngestWithResult validates, stores, and returns counts.
func (ing *CurriculumIngester) IngestWithResult(ctx context.Context, rawJSON json.RawMessage) (*IngestResult, error) {
	if err := schema.Validate(rawJSON); err != nil {
		return nil, fmt.Errorf("schema validation: %w", err)
	}

	var curr CurriculumOutput
	if err := json.Unmarshal(rawJSON, &curr); err != nil {
		return nil, fmt.Errorf("unmarshal curriculum: %w", err)
	}

	result := &IngestResult{}

	tx, err := ing.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	if err := ing.storeTopic(ctx, tx, &curr); err != nil {
		return nil, err
	}

	for i, mod := range curr.Modules {
		if err := ing.storeModule(ctx, tx, curr.ID, &mod, i); err != nil {
			return nil, err
		}

		result.ModulesCreated++

		for j, lesson := range mod.Lessons {
			if err := ing.storeLesson(ctx, tx, mod.ID, &lesson, j); err != nil {
				return nil, err
			}

			result.LessonsCreated++

			for _, concept := range lesson.ConceptsTaught {
				if err := ing.storeConcept(ctx, tx, curr.ID, lesson.ID, &concept); err != nil {
					return nil, err
				}

				result.ConceptsCreated++
			}

			for _, ref := range lesson.ConceptsReferenced {
				if err := ing.storeConceptReference(ctx, tx, ref.ID, lesson.ID); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := ing.storePrerequisites(ctx, tx, curr.ID, &curr.Prerequisites); err != nil {
		return nil, err
	}

	if err := ing.storeExpansionQueue(ctx, tx, curr.ID, &curr.Prerequisites); err != nil {
		return nil, err
	}

	if err := ing.storeSearchIndex(ctx, tx, &curr); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return result, nil
}
