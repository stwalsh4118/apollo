package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sean/apollo/api/internal/models"
)

// GraphRepository defines graph data operations.
type GraphRepository interface {
	GetFullGraph(ctx context.Context) (*models.GraphData, error)
	GetTopicGraph(ctx context.Context, topicID string) (*models.GraphData, error)
}

// SQLiteGraphRepository implements GraphRepository using SQLite.
type SQLiteGraphRepository struct {
	db *sql.DB
}

// NewGraphRepository creates a new SQLiteGraphRepository.
func NewGraphRepository(db *sql.DB) *SQLiteGraphRepository {
	return &SQLiteGraphRepository{db: db}
}

const topicNodesSQL = `SELECT id, title, 'topic' FROM topics`
const conceptNodesSQL = `SELECT id, name, 'concept' FROM concepts`

const prerequisiteEdgesSQL = `
SELECT topic_id, prerequisite_topic_id, priority
FROM topic_prerequisites
`

// referenceEdgesSQL creates edges from concepts to their defining topic, but only
// for concepts that are actually referenced in at least one lesson (via concept_references).
// This keeps the graph focused on actively-used concepts rather than all defined concepts.
const referenceEdgesSQL = `
SELECT c.id, c.defined_in_topic, 'reference'
FROM concept_references cr
JOIN concepts c ON c.id = cr.concept_id
WHERE c.defined_in_topic IS NOT NULL
GROUP BY c.id, c.defined_in_topic
`

const relationEdgesSQL = `
SELECT topic_a, topic_b, relation_type
FROM topic_relations
`

func (r *SQLiteGraphRepository) GetFullGraph(ctx context.Context) (*models.GraphData, error) {
	nodes, err := r.queryNodes(ctx, topicNodesSQL, conceptNodesSQL)
	if err != nil {
		return nil, err
	}

	edges, err := r.queryAllEdges(ctx)
	if err != nil {
		return nil, err
	}

	return &models.GraphData{Nodes: nodes, Edges: edges}, nil
}

const topicExistsSQL = `SELECT EXISTS(SELECT 1 FROM topics WHERE id = ?)`

const topicSubgraphNodesSQL = `
SELECT id, title, 'topic' FROM topics WHERE id = ?
UNION ALL
SELECT id, name, 'concept' FROM concepts WHERE defined_in_topic = ?
`

const topicSubgraphPrereqSQL = `
SELECT topic_id, prerequisite_topic_id, priority
FROM topic_prerequisites
WHERE topic_id = ? OR prerequisite_topic_id = ?
`

const topicSubgraphRelationSQL = `
SELECT topic_a, topic_b, relation_type
FROM topic_relations
WHERE topic_a = ? OR topic_b = ?
`

const topicSubgraphRefSQL = `
SELECT c.id, c.defined_in_topic, 'reference'
FROM concept_references cr
JOIN concepts c ON c.id = cr.concept_id
WHERE c.defined_in_topic = ?
GROUP BY c.id, c.defined_in_topic
`

func (r *SQLiteGraphRepository) GetTopicGraph(ctx context.Context, topicID string) (*models.GraphData, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, topicExistsSQL, topicID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check topic exists: %w", err)
	}

	if !exists {
		return nil, nil
	}

	nodes, err := r.queryNodeRows(ctx, topicSubgraphNodesSQL, topicID, topicID)
	if err != nil {
		return nil, fmt.Errorf("query topic subgraph nodes: %w", err)
	}

	var edges []models.GraphEdge

	prereqEdges, err := r.queryEdgeRows(ctx, topicSubgraphPrereqSQL, topicID, topicID)
	if err != nil {
		return nil, err
	}

	edges = append(edges, prereqEdges...)

	relEdges, err := r.queryEdgeRows(ctx, topicSubgraphRelationSQL, topicID, topicID)
	if err != nil {
		return nil, err
	}

	edges = append(edges, relEdges...)

	refEdges, err := r.queryEdgeRows(ctx, topicSubgraphRefSQL, topicID)
	if err != nil {
		return nil, err
	}

	edges = append(edges, refEdges...)

	if edges == nil {
		edges = []models.GraphEdge{}
	}

	return &models.GraphData{Nodes: nodes, Edges: edges}, nil
}

func (r *SQLiteGraphRepository) queryNodes(ctx context.Context, topicSQL, conceptSQL string) ([]models.GraphNode, error) {
	topicNodes, err := r.queryNodeRows(ctx, topicSQL)
	if err != nil {
		return nil, fmt.Errorf("query topic nodes: %w", err)
	}

	conceptNodes, err := r.queryNodeRows(ctx, conceptSQL)
	if err != nil {
		return nil, fmt.Errorf("query concept nodes: %w", err)
	}

	nodes := append(topicNodes, conceptNodes...)
	if nodes == nil {
		nodes = []models.GraphNode{}
	}

	return nodes, nil
}

func (r *SQLiteGraphRepository) queryNodeRows(ctx context.Context, query string, args ...any) ([]models.GraphNode, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query nodes: %w", err)
	}
	defer rows.Close()

	var nodes []models.GraphNode

	for rows.Next() {
		var n models.GraphNode
		if err := rows.Scan(&n.ID, &n.Label, &n.Type); err != nil {
			return nil, fmt.Errorf("scan node: %w", err)
		}

		nodes = append(nodes, n)
	}

	return nodes, rows.Err()
}

func (r *SQLiteGraphRepository) queryAllEdges(ctx context.Context) ([]models.GraphEdge, error) {
	var edges []models.GraphEdge

	prereqs, err := r.queryEdgeRows(ctx, prerequisiteEdgesSQL)
	if err != nil {
		return nil, fmt.Errorf("query prerequisite edges: %w", err)
	}

	edges = append(edges, prereqs...)

	refs, err := r.queryEdgeRows(ctx, referenceEdgesSQL)
	if err != nil {
		return nil, fmt.Errorf("query reference edges: %w", err)
	}

	edges = append(edges, refs...)

	rels, err := r.queryEdgeRows(ctx, relationEdgesSQL)
	if err != nil {
		return nil, fmt.Errorf("query relation edges: %w", err)
	}

	edges = append(edges, rels...)

	if edges == nil {
		edges = []models.GraphEdge{}
	}

	return edges, nil
}

func (r *SQLiteGraphRepository) queryEdgeRows(ctx context.Context, query string, args ...any) ([]models.GraphEdge, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query edges: %w", err)
	}
	defer rows.Close()

	var edges []models.GraphEdge

	for rows.Next() {
		var e models.GraphEdge
		if err := rows.Scan(&e.Source, &e.Target, &e.Type); err != nil {
			return nil, fmt.Errorf("scan edge: %w", err)
		}

		edges = append(edges, e)
	}

	return edges, rows.Err()
}

// Verify interface compliance at compile time.
var _ GraphRepository = (*SQLiteGraphRepository)(nil)
