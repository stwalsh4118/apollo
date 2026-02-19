package schema

import (
	"strings"
	"testing"
)

// validMinimalCurriculum is the smallest valid Topic JSON that satisfies the schema.
const validMinimalCurriculum = `{
  "id": "test-topic",
  "title": "Test Topic",
  "description": "A minimal test topic.",
  "difficulty": "foundational",
  "estimated_hours": 1,
  "tags": ["test"],
  "prerequisites": {
    "essential": [],
    "helpful": [],
    "deep_background": []
  },
  "related_topics": [],
  "modules": [
    {
      "id": "test-topic/intro",
      "title": "Introduction",
      "description": "Intro module.",
      "learning_objectives": ["Understand basics"],
      "estimated_minutes": 30,
      "order": 1,
      "lessons": [
        {
          "id": "test-topic/intro/hello",
          "title": "Hello",
          "order": 1,
          "estimated_minutes": 10,
          "content": {
            "sections": [
              { "type": "text", "body": "Welcome to the topic." }
            ]
          },
          "concepts_taught": [
            {
              "id": "test-concept",
              "name": "Test Concept",
              "definition": "A concept for testing.",
              "flashcard": {
                "front": "What is a test concept?",
                "back": "A concept used for testing."
              }
            }
          ],
          "concepts_referenced": [],
          "examples": [],
          "exercises": [],
          "review_questions": []
        }
      ],
      "assessment": {
        "questions": [
          {
            "type": "conceptual",
            "question": "What is this topic about?",
            "answer": "Testing.",
            "concepts_tested": ["test-concept"]
          }
        ]
      }
    }
  ],
  "source_urls": ["https://example.com"],
  "generated_at": "2026-02-19T00:00:00Z",
  "version": 1
}`

func TestValidate_ValidCurriculum(t *testing.T) {
	err := Validate([]byte(validMinimalCurriculum))
	if err != nil {
		t.Fatalf("expected valid curriculum to pass, got: %v", err)
	}
}

func TestValidate_MissingRequiredID(t *testing.T) {
	// Topic missing "id" field.
	input := `{
		"title": "No ID",
		"description": "Missing id.",
		"difficulty": "foundational",
		"estimated_hours": 1,
		"tags": [],
		"prerequisites": { "essential": [], "helpful": [], "deep_background": [] },
		"related_topics": [],
		"modules": [
			{
				"id": "m", "title": "M", "description": "D",
				"learning_objectives": ["L"], "estimated_minutes": 10, "order": 1,
				"lessons": [{
					"id": "l", "title": "L", "order": 1, "estimated_minutes": 5,
					"content": { "sections": [{ "type": "text", "body": "B" }] },
					"concepts_taught": [], "concepts_referenced": [],
					"examples": [], "exercises": [], "review_questions": []
				}],
				"assessment": { "questions": [{ "type": "conceptual", "question": "Q", "answer": "A", "concepts_tested": [] }] }
			}
		],
		"source_urls": [],
		"generated_at": "2026-02-19T00:00:00Z",
		"version": 1
	}`

	err := Validate([]byte(input))
	if err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
	if !strings.Contains(err.Error(), "id") {
		t.Errorf("expected error to mention 'id', got: %v", err)
	}
}

func TestValidate_InvalidDifficulty(t *testing.T) {
	// Topic with invalid difficulty enum.
	input := `{
		"id": "t", "title": "T", "description": "D",
		"difficulty": "unknown_level",
		"estimated_hours": 1,
		"tags": [],
		"prerequisites": { "essential": [], "helpful": [], "deep_background": [] },
		"related_topics": [],
		"modules": [
			{
				"id": "m", "title": "M", "description": "D",
				"learning_objectives": ["L"], "estimated_minutes": 10, "order": 1,
				"lessons": [{
					"id": "l", "title": "L", "order": 1, "estimated_minutes": 5,
					"content": { "sections": [{ "type": "text", "body": "B" }] },
					"concepts_taught": [], "concepts_referenced": [],
					"examples": [], "exercises": [], "review_questions": []
				}],
				"assessment": { "questions": [{ "type": "conceptual", "question": "Q", "answer": "A", "concepts_tested": [] }] }
			}
		],
		"source_urls": [],
		"generated_at": "2026-02-19T00:00:00Z",
		"version": 1
	}`

	err := Validate([]byte(input))
	if err == nil {
		t.Fatal("expected error for invalid difficulty, got nil")
	}
}

func TestValidate_WrongType(t *testing.T) {
	// estimated_hours as string instead of number.
	input := `{
		"id": "t", "title": "T", "description": "D",
		"difficulty": "foundational",
		"estimated_hours": "not_a_number",
		"tags": [],
		"prerequisites": { "essential": [], "helpful": [], "deep_background": [] },
		"related_topics": [],
		"modules": [
			{
				"id": "m", "title": "M", "description": "D",
				"learning_objectives": ["L"], "estimated_minutes": 10, "order": 1,
				"lessons": [{
					"id": "l", "title": "L", "order": 1, "estimated_minutes": 5,
					"content": { "sections": [{ "type": "text", "body": "B" }] },
					"concepts_taught": [], "concepts_referenced": [],
					"examples": [], "exercises": [], "review_questions": []
				}],
				"assessment": { "questions": [{ "type": "conceptual", "question": "Q", "answer": "A", "concepts_tested": [] }] }
			}
		],
		"source_urls": [],
		"generated_at": "2026-02-19T00:00:00Z",
		"version": 1
	}`

	err := Validate([]byte(input))
	if err == nil {
		t.Fatal("expected error for wrong type, got nil")
	}
}

func TestValidate_InvalidJSON(t *testing.T) {
	err := Validate([]byte(`{not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("expected 'invalid JSON' in error, got: %v", err)
	}
}

func TestValidate_ErrorContainsFieldPath(t *testing.T) {
	// Missing required 'title' in a module — error should reference the path.
	input := `{
		"id": "t", "title": "T", "description": "D",
		"difficulty": "foundational",
		"estimated_hours": 1,
		"tags": [],
		"prerequisites": { "essential": [], "helpful": [], "deep_background": [] },
		"related_topics": [],
		"modules": [
			{
				"id": "m", "description": "D",
				"learning_objectives": ["L"], "estimated_minutes": 10, "order": 1,
				"lessons": [{
					"id": "l", "title": "L", "order": 1, "estimated_minutes": 5,
					"content": { "sections": [{ "type": "text", "body": "B" }] },
					"concepts_taught": [], "concepts_referenced": [],
					"examples": [], "exercises": [], "review_questions": []
				}],
				"assessment": { "questions": [{ "type": "conceptual", "question": "Q", "answer": "A", "concepts_tested": [] }] }
			}
		],
		"source_urls": [],
		"generated_at": "2026-02-19T00:00:00Z",
		"version": 1
	}`

	err := Validate([]byte(input))
	if err == nil {
		t.Fatal("expected error for missing module title, got nil")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "modules") {
		t.Errorf("expected error to contain field path with 'modules', got: %v", errStr)
	}
}

func TestValidate_SchemaLoadedOnce(t *testing.T) {
	// Call Validate twice — schema should be cached (sync.Once).
	err1 := Validate([]byte(validMinimalCurriculum))
	err2 := Validate([]byte(validMinimalCurriculum))
	if err1 != nil {
		t.Fatalf("first call failed: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("second call failed: %v", err2)
	}
}
