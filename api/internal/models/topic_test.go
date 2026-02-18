package models_test

import (
	"encoding/json"
	"testing"

	"github.com/sean/apollo/api/internal/models"
)

func TestParseJSONStringSlice(t *testing.T) {
	tests := []struct {
		name string
		raw  *string
		want int
	}{
		{"nil", nil, 0},
		{"empty string", strPtr(""), 0},
		{"null string", strPtr("null"), 0},
		{"valid array", strPtr(`["go","python","rust"]`), 3},
		{"invalid json", strPtr("not json"), 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := models.ParseJSONStringSlice(tc.raw)

			if len(result) != tc.want {
				t.Fatalf("expected %d items, got %d", tc.want, len(result))
			}
		})
	}
}

func TestTopicDetailEmbeddedFieldsSerialization(t *testing.T) {
	td := models.TopicDetail{
		Modules: []models.ModuleSummary{{ID: "m-1", Title: "Module 1", SortOrder: 1}},
	}
	td.ID = "test-topic"
	td.Title = "Test Topic"
	td.Status = "published"

	data, err := json.Marshal(td)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// Verify embedded fields flatten correctly.
	if m["id"] != "test-topic" {
		t.Fatalf("expected embedded id 'test-topic', got %v", m["id"])
	}

	modules, ok := m["modules"].([]any)
	if !ok || len(modules) != 1 {
		t.Fatalf("expected 1 module, got %v", m["modules"])
	}
}

func TestTopicSummaryJSONShape(t *testing.T) {
	ts := models.TopicSummary{
		ID:          "golang-basics",
		Title:       "Go Basics",
		Difficulty:  "foundational",
		Status:      "published",
		ModuleCount: 5,
		Tags:        []string{"go", "programming"},
	}

	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if m["id"] != "golang-basics" {
		t.Fatalf("expected id 'golang-basics', got %v", m["id"])
	}

	if m["module_count"].(float64) != 5 {
		t.Fatalf("expected module_count 5, got %v", m["module_count"])
	}
}

func strPtr(s string) *string {
	return &s
}
