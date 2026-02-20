package research

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sean/apollo/api/internal/schema"
)

const poolSummaryFilename = "knowledge_pool_summary.json"

const listTopicIDsSQL = `SELECT id FROM topics ORDER BY id`

const listModuleIDsForTopicSQL = `SELECT id FROM modules WHERE topic_id = ? ORDER BY sort_order`

const listConceptIDsSQL = `SELECT id FROM concepts ORDER BY id`

// PoolSummaryTopic is one entry in existing_topics.
type PoolSummaryTopic struct {
	ID      string   `json:"id"`
	Modules []string `json:"modules"`
}

// PoolSummary is the JSON structure written to knowledge_pool_summary.json.
type PoolSummary struct {
	ExistingTopics   []PoolSummaryTopic `json:"existing_topics"`
	ExistingConcepts []string           `json:"existing_concepts"`
}

// PoolSummaryBuilder queries the database and produces the knowledge pool summary.
type PoolSummaryBuilder struct {
	db *sql.DB
}

// NewPoolSummaryBuilder creates a new PoolSummaryBuilder.
func NewPoolSummaryBuilder(db *sql.DB) *PoolSummaryBuilder {
	return &PoolSummaryBuilder{db: db}
}

// Build queries the database and returns the pool summary as JSON bytes.
func (b *PoolSummaryBuilder) Build(ctx context.Context) ([]byte, error) {
	topics, err := b.queryTopics(ctx)
	if err != nil {
		return nil, err
	}

	concepts, err := b.queryConcepts(ctx)
	if err != nil {
		return nil, err
	}

	summary := PoolSummary{
		ExistingTopics:   topics,
		ExistingConcepts: concepts,
	}

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal pool summary: %w", err)
	}

	if err := schema.ValidatePoolSummary(data); err != nil {
		return nil, fmt.Errorf("validate pool summary: %w", err)
	}

	return data, nil
}

// WriteToDir builds the summary and writes it to the specified directory.
func (b *PoolSummaryBuilder) WriteToDir(ctx context.Context, dir string) error {
	data, err := b.Build(ctx)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, poolSummaryFilename)

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write pool summary to %s: %w", path, err)
	}

	return nil
}

func (b *PoolSummaryBuilder) queryTopics(ctx context.Context) ([]PoolSummaryTopic, error) {
	rows, err := b.db.QueryContext(ctx, listTopicIDsSQL)
	if err != nil {
		return nil, fmt.Errorf("query topic IDs: %w", err)
	}
	defer rows.Close()

	var topicIDs []string

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan topic ID: %w", err)
		}

		topicIDs = append(topicIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate topic IDs: %w", err)
	}

	topics := make([]PoolSummaryTopic, 0, len(topicIDs))

	for _, topicID := range topicIDs {
		modules, err := b.queryModules(ctx, topicID)
		if err != nil {
			return nil, err
		}

		topics = append(topics, PoolSummaryTopic{
			ID:      topicID,
			Modules: modules,
		})
	}

	return topics, nil
}

func (b *PoolSummaryBuilder) queryModules(ctx context.Context, topicID string) ([]string, error) {
	rows, err := b.db.QueryContext(ctx, listModuleIDsForTopicSQL, topicID)
	if err != nil {
		return nil, fmt.Errorf("query module IDs for topic %s: %w", topicID, err)
	}
	defer rows.Close()

	var modules []string

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan module ID: %w", err)
		}

		modules = append(modules, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate module IDs: %w", err)
	}

	if modules == nil {
		modules = []string{}
	}

	return modules, nil
}

func (b *PoolSummaryBuilder) queryConcepts(ctx context.Context) ([]string, error) {
	rows, err := b.db.QueryContext(ctx, listConceptIDsSQL)
	if err != nil {
		return nil, fmt.Errorf("query concept IDs: %w", err)
	}
	defer rows.Close()

	var concepts []string

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan concept ID: %w", err)
		}

		concepts = append(concepts, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate concept IDs: %w", err)
	}

	if concepts == nil {
		concepts = []string{}
	}

	return concepts, nil
}
