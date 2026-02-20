package research_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sean/apollo/api/internal/research"
)

// sampleCurriculum is a valid curriculum JSON matching the schema.
const sampleCurriculum = `{
  "id": "go-concurrency",
  "title": "Go Concurrency",
  "description": "Learn concurrent programming in Go.",
  "difficulty": "intermediate",
  "estimated_hours": 10,
  "tags": ["go", "concurrency"],
  "prerequisites": {
    "essential": [
      {"topic_id": "go-basics", "reason": "Need Go fundamentals"}
    ],
    "helpful": [
      {"topic_id": "os-threads", "reason": "Understanding OS threads helps"}
    ],
    "deep_background": [
      {"topic_id": "csp-theory", "reason": "CSP theory background"}
    ]
  },
  "related_topics": ["go-networking"],
  "modules": [
    {
      "id": "go-concurrency/goroutines",
      "title": "Goroutines",
      "description": "Lightweight threads in Go.",
      "learning_objectives": ["Understand goroutines"],
      "estimated_minutes": 60,
      "order": 1,
      "lessons": [
        {
          "id": "go-concurrency/goroutines/intro",
          "title": "Introduction to Goroutines",
          "order": 1,
          "estimated_minutes": 30,
          "content": {"sections": [{"type": "text", "body": "Goroutines are lightweight threads."}]},
          "concepts_taught": [
            {
              "id": "goroutine",
              "name": "Goroutine",
              "definition": "A lightweight thread managed by the Go runtime.",
              "flashcard": {
                "front": "What is a goroutine?",
                "back": "A lightweight thread managed by the Go runtime."
              }
            }
          ],
          "concepts_referenced": [],
          "examples": [{"title": "Basic goroutine", "description": "Launch a goroutine", "code": "go func() {}()", "explanation": "Creates a new goroutine."}],
          "exercises": [{"type": "command", "title": "Run goroutine", "instructions": "Launch a goroutine", "success_criteria": ["Goroutine runs"], "hints": ["Use go keyword"], "environment": "terminal"}],
          "review_questions": [{"question": "What is a goroutine?", "answer": "A lightweight thread.", "concepts_tested": ["goroutine"]}]
        },
        {
          "id": "go-concurrency/goroutines/sync",
          "title": "Synchronizing Goroutines",
          "order": 2,
          "estimated_minutes": 30,
          "content": {"sections": [{"type": "text", "body": "Use WaitGroup for synchronization."}]},
          "concepts_taught": [
            {
              "id": "waitgroup",
              "name": "WaitGroup",
              "definition": "A synchronization primitive for waiting on goroutines.",
              "flashcard": {
                "front": "What is sync.WaitGroup?",
                "back": "A synchronization primitive for waiting on goroutines."
              }
            }
          ],
          "concepts_referenced": [
            {"id": "goroutine", "defined_in": "go-concurrency/goroutines/intro"}
          ],
          "examples": [],
          "exercises": [],
          "review_questions": []
        }
      ],
      "assessment": {"questions": [{"type": "conceptual", "question": "Explain goroutines.", "answer": "Lightweight threads.", "concepts_tested": ["goroutine"]}]}
    }
  ],
  "source_urls": ["https://go.dev/doc"],
  "generated_at": "2026-02-19T08:00:00Z",
  "version": 1
}`

func TestIngestValidCurriculum(t *testing.T) {
	db := setupTestDB(t)
	ingester := research.NewCurriculumIngester(db)

	if err := ingester.Ingest(context.Background(), json.RawMessage(sampleCurriculum)); err != nil {
		t.Fatalf("ingest: %v", err)
	}

	// Verify topic created.
	var topicTitle string
	if err := db.QueryRow("SELECT title FROM topics WHERE id = ?", "go-concurrency").Scan(&topicTitle); err != nil {
		t.Fatalf("query topic: %v", err)
	}

	if topicTitle != "Go Concurrency" {
		t.Fatalf("expected title 'Go Concurrency', got %q", topicTitle)
	}

	// Verify module created.
	var moduleCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM modules WHERE topic_id = ?", "go-concurrency").Scan(&moduleCount); err != nil {
		t.Fatalf("count modules: %v", err)
	}

	if moduleCount != 1 {
		t.Fatalf("expected 1 module, got %d", moduleCount)
	}

	// Verify lessons created.
	var lessonCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM lessons l JOIN modules m ON l.module_id = m.id WHERE m.topic_id = ?", "go-concurrency").Scan(&lessonCount); err != nil {
		t.Fatalf("count lessons: %v", err)
	}

	if lessonCount != 2 {
		t.Fatalf("expected 2 lessons, got %d", lessonCount)
	}

	// Verify concepts created.
	var conceptCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM concepts WHERE defined_in_topic = ?", "go-concurrency").Scan(&conceptCount); err != nil {
		t.Fatalf("count concepts: %v", err)
	}

	if conceptCount != 2 {
		t.Fatalf("expected 2 concepts, got %d", conceptCount)
	}

	// Verify concept has flashcard.
	var front, back string
	if err := db.QueryRow("SELECT flashcard_front, flashcard_back FROM concepts WHERE id = ?", "goroutine").Scan(&front, &back); err != nil {
		t.Fatalf("query concept flashcard: %v", err)
	}

	if front != "What is a goroutine?" {
		t.Fatalf("expected flashcard front 'What is a goroutine?', got %q", front)
	}

	// Verify concept references (goroutine appears in both lessons).
	var refCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM concept_references WHERE concept_id = ?", "goroutine").Scan(&refCount); err != nil {
		t.Fatalf("count concept refs: %v", err)
	}

	if refCount != 2 {
		t.Fatalf("expected 2 references for goroutine (defining + referencing), got %d", refCount)
	}

	// Verify prerequisites stored (only those whose target topic exists; none exist in test DB).
	var prereqCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM topic_prerequisites WHERE topic_id = ?", "go-concurrency").Scan(&prereqCount); err != nil {
		t.Fatalf("count prerequisites: %v", err)
	}

	if prereqCount != 0 {
		t.Fatalf("expected 0 prerequisites (target topics don't exist), got %d", prereqCount)
	}

	// Verify expansion queue (helpful + deep_background only).
	var expansionCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM expansion_queue WHERE requested_by_topic = ?", "go-concurrency").Scan(&expansionCount); err != nil {
		t.Fatalf("count expansion queue: %v", err)
	}

	if expansionCount != 2 {
		t.Fatalf("expected 2 expansion queue entries (helpful + deep_background), got %d", expansionCount)
	}

	// Verify expansion queue statuses are 'available'.
	var availableCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM expansion_queue WHERE status = 'available'").Scan(&availableCount); err != nil {
		t.Fatalf("count available: %v", err)
	}

	if availableCount != 2 {
		t.Fatalf("expected 2 available entries, got %d", availableCount)
	}

	// Verify search index updated.
	var searchCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM search_index WHERE entity_id = ?", "go-concurrency").Scan(&searchCount); err != nil {
		t.Fatalf("count search index: %v", err)
	}

	if searchCount != 1 {
		t.Fatalf("expected 1 search entry for topic, got %d", searchCount)
	}
}

func TestIngestInvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	ingester := research.NewCurriculumIngester(db)

	err := ingester.Ingest(context.Background(), json.RawMessage(`{"bad": "data"}`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	// Verify nothing was inserted.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM topics").Scan(&count); err != nil {
		t.Fatalf("count topics: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected 0 topics after failed validation, got %d", count)
	}
}

func TestIngestTransactionRollback(t *testing.T) {
	db := setupTestDB(t)
	ingester := research.NewCurriculumIngester(db)

	// First ingest succeeds.
	if err := ingester.Ingest(context.Background(), json.RawMessage(sampleCurriculum)); err != nil {
		t.Fatalf("first ingest: %v", err)
	}

	// Second ingest with same data fails (duplicate topic ID).
	err := ingester.Ingest(context.Background(), json.RawMessage(sampleCurriculum))
	if err == nil {
		t.Fatal("expected error for duplicate ingest")
	}

	// Verify only one topic exists (rollback worked).
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM topics").Scan(&count); err != nil {
		t.Fatalf("count topics: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 topic after rollback, got %d", count)
	}
}

func TestIngestWithResult(t *testing.T) {
	db := setupTestDB(t)
	ingester := research.NewCurriculumIngester(db)

	result, err := ingester.IngestWithResult(context.Background(), json.RawMessage(sampleCurriculum))
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}

	if result.ModulesCreated != 1 {
		t.Fatalf("expected 1 module, got %d", result.ModulesCreated)
	}

	if result.LessonsCreated != 2 {
		t.Fatalf("expected 2 lessons, got %d", result.LessonsCreated)
	}

	if result.ConceptsCreated != 2 {
		t.Fatalf("expected 2 concepts, got %d", result.ConceptsCreated)
	}
}

func TestIngestModuleSortOrder(t *testing.T) {
	db := setupTestDB(t)
	ingester := research.NewCurriculumIngester(db)

	if err := ingester.Ingest(context.Background(), json.RawMessage(sampleCurriculum)); err != nil {
		t.Fatalf("ingest: %v", err)
	}

	var sortOrder int
	if err := db.QueryRow("SELECT sort_order FROM modules WHERE id = ?", "go-concurrency/goroutines").Scan(&sortOrder); err != nil {
		t.Fatalf("query sort_order: %v", err)
	}

	if sortOrder != 1 {
		t.Fatalf("expected sort_order 1, got %d", sortOrder)
	}

	// Verify lesson sort orders.
	rows, err := db.Query("SELECT id, sort_order FROM lessons ORDER BY sort_order")
	if err != nil {
		t.Fatalf("query lessons: %v", err)
	}
	defer rows.Close()

	expected := map[string]int{
		"go-concurrency/goroutines/intro": 1,
		"go-concurrency/goroutines/sync":  2,
	}

	for rows.Next() {
		var id string
		var so int
		if err := rows.Scan(&id, &so); err != nil {
			t.Fatalf("scan: %v", err)
		}

		if exp, ok := expected[id]; ok && exp != so {
			t.Fatalf("lesson %s: expected sort_order %d, got %d", id, exp, so)
		}
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("rows err: %v", err)
	}
}
