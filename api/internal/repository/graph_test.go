package repository_test

import (
	"context"
	"testing"

	"github.com/sean/apollo/api/internal/repository"
)

func TestGetFullGraphNodes(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGraphRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedTopic(t, db, "t2", "Advanced Go", "advanced", "published")
	seedConcept(t, db, "c1", "Variable", "Storage location", "t1")
	seedConcept(t, db, "c2", "Function", "Reusable code", "t1")
	seedConcept(t, db, "c3", "Goroutine", "Lightweight thread", "t2")

	graph, err := repo.GetFullGraph(context.Background())
	if err != nil {
		t.Fatalf("get full graph: %v", err)
	}

	if len(graph.Nodes) != 5 {
		t.Fatalf("expected 5 nodes (2 topics + 3 concepts), got %d", len(graph.Nodes))
	}

	topicCount := 0
	conceptCount := 0

	for _, n := range graph.Nodes {
		switch n.Type {
		case "topic":
			topicCount++
		case "concept":
			conceptCount++
		}
	}

	if topicCount != 2 {
		t.Fatalf("expected 2 topic nodes, got %d", topicCount)
	}

	if conceptCount != 3 {
		t.Fatalf("expected 3 concept nodes, got %d", conceptCount)
	}
}

func TestGetFullGraphEdges(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGraphRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedTopic(t, db, "t2", "Advanced Go", "advanced", "published")
	seedModule(t, db, "m1", "t1", "Intro", 1)
	seedLesson(t, db, "l1", "m1", "Hello World", 1)
	seedConcept(t, db, "c1", "Variable", "Storage location", "t1")
	seedConceptReference(t, db, "c1", "l1")

	// Prerequisite: t2 requires t1.
	mustExec(t, db,
		"INSERT INTO topic_prerequisites (topic_id, prerequisite_topic_id, priority) VALUES (?, ?, ?)",
		"t2", "t1", "essential",
	)

	// Relation: t1 related to t2.
	mustExec(t, db,
		"INSERT INTO topic_relations (topic_a, topic_b, relation_type) VALUES (?, ?, ?)",
		"t1", "t2", "builds_on",
	)

	graph, err := repo.GetFullGraph(context.Background())
	if err != nil {
		t.Fatalf("get full graph: %v", err)
	}

	if len(graph.Edges) < 3 {
		t.Fatalf("expected at least 3 edges (prereq + ref + relation), got %d", len(graph.Edges))
	}

	edgeTypes := make(map[string]bool)
	for _, e := range graph.Edges {
		edgeTypes[e.Type] = true
	}

	if !edgeTypes["essential"] {
		t.Fatal("expected prerequisite edge with type 'essential'")
	}

	if !edgeTypes["reference"] {
		t.Fatal("expected reference edge with type 'reference'")
	}

	if !edgeTypes["builds_on"] {
		t.Fatal("expected relation edge with type 'builds_on'")
	}
}

func TestGetTopicGraph(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGraphRepository(db)

	seedTopic(t, db, "t1", "Go Basics", "foundational", "published")
	seedTopic(t, db, "t2", "Advanced Go", "advanced", "published")
	seedConcept(t, db, "c1", "Variable", "Storage location", "t1")
	seedConcept(t, db, "c2", "Function", "Reusable code", "t1")
	seedConcept(t, db, "c3", "Goroutine", "Lightweight thread", "t2")

	graph, err := repo.GetTopicGraph(context.Background(), "t1")
	if err != nil {
		t.Fatalf("get topic graph: %v", err)
	}

	if graph == nil {
		t.Fatal("expected non-nil graph")
	}

	// 1 topic node + 2 concepts for t1.
	if len(graph.Nodes) != 3 {
		t.Fatalf("expected 3 nodes (1 topic + 2 concepts), got %d", len(graph.Nodes))
	}
}

func TestGetTopicGraphNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGraphRepository(db)

	graph, err := repo.GetTopicGraph(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if graph != nil {
		t.Fatal("expected nil graph for nonexistent topic")
	}
}

func TestGetFullGraphEmpty(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGraphRepository(db)

	graph, err := repo.GetFullGraph(context.Background())
	if err != nil {
		t.Fatalf("get full graph: %v", err)
	}

	if len(graph.Nodes) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.Edges) != 0 {
		t.Fatalf("expected 0 edges, got %d", len(graph.Edges))
	}
}
