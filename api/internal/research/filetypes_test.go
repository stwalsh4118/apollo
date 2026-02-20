package research

import (
	"encoding/json"
	"testing"
)

func TestTopicFileJSONRoundTrip(t *testing.T) {
	original := TopicFile{
		ID:             "go-concurrency",
		Title:          "Go Concurrency",
		Description:    "Learn concurrent programming in Go",
		Difficulty:     "intermediate",
		EstimatedHours: 12.5,
		Tags:           []string{"go", "concurrency"},
		Prerequisites: PrerequisitesOutput{
			Essential:      []PrerequisiteItem{{TopicID: "go-basics", Reason: "core language"}},
			Helpful:        []PrerequisiteItem{{TopicID: "os-fundamentals", Reason: "system calls"}},
			DeepBackground: nil,
		},
		RelatedTopics: []string{"distributed-systems"},
		SourceURLs:    []string{"https://go.dev/doc"},
		GeneratedAt:   "2026-02-20T10:00:00Z",
		Version:       1,
		ModulePlan: []ModulePlanEntry{
			{ID: "goroutines", Title: "Goroutines", Description: "Lightweight threads", Order: 1},
			{ID: "channels", Title: "Channels", Description: "Communication between goroutines", Order: 2},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal TopicFile: %v", err)
	}

	var decoded TopicFile
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal TopicFile: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %q, want %q", decoded.ID, original.ID)
	}
	if decoded.Title != original.Title {
		t.Errorf("Title = %q, want %q", decoded.Title, original.Title)
	}
	if decoded.Difficulty != original.Difficulty {
		t.Errorf("Difficulty = %q, want %q", decoded.Difficulty, original.Difficulty)
	}
	if decoded.EstimatedHours != original.EstimatedHours {
		t.Errorf("EstimatedHours = %v, want %v", decoded.EstimatedHours, original.EstimatedHours)
	}
	if len(decoded.ModulePlan) != len(original.ModulePlan) {
		t.Fatalf("ModulePlan length = %d, want %d", len(decoded.ModulePlan), len(original.ModulePlan))
	}
	if decoded.ModulePlan[0].ID != "goroutines" {
		t.Errorf("ModulePlan[0].ID = %q, want %q", decoded.ModulePlan[0].ID, "goroutines")
	}
	if decoded.ModulePlan[1].Order != 2 {
		t.Errorf("ModulePlan[1].Order = %d, want %d", decoded.ModulePlan[1].Order, 2)
	}
	if len(decoded.Prerequisites.Essential) != 1 {
		t.Errorf("Prerequisites.Essential length = %d, want 1", len(decoded.Prerequisites.Essential))
	}
}

func TestModulePlanEntryJSONRoundTrip(t *testing.T) {
	original := ModulePlanEntry{
		ID:          "goroutines",
		Title:       "Goroutines",
		Description: "Lightweight threads in Go",
		Order:       1,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal ModulePlanEntry: %v", err)
	}

	var decoded ModulePlanEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal ModulePlanEntry: %v", err)
	}

	if decoded != original {
		t.Errorf("decoded = %+v, want %+v", decoded, original)
	}
}

func TestModuleFileJSONRoundTrip(t *testing.T) {
	original := ModuleFile{
		ID:                 "goroutines",
		Title:              "Goroutines",
		Description:        "Lightweight threads in Go",
		Order:              1,
		LearningObjectives: []string{"Understand goroutine lifecycle", "Use go keyword"},
		EstimatedMinutes:   90,
		Assessment:         json.RawMessage(`{"questions":[{"q":"What is a goroutine?"}]}`),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal ModuleFile: %v", err)
	}

	var decoded ModuleFile
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal ModuleFile: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %q, want %q", decoded.ID, original.ID)
	}
	if decoded.Title != original.Title {
		t.Errorf("Title = %q, want %q", decoded.Title, original.Title)
	}
	if decoded.Order != original.Order {
		t.Errorf("Order = %d, want %d", decoded.Order, original.Order)
	}
	if len(decoded.LearningObjectives) != 2 {
		t.Errorf("LearningObjectives length = %d, want 2", len(decoded.LearningObjectives))
	}
	if decoded.EstimatedMinutes != original.EstimatedMinutes {
		t.Errorf("EstimatedMinutes = %d, want %d", decoded.EstimatedMinutes, original.EstimatedMinutes)
	}
	if string(decoded.Assessment) != string(original.Assessment) {
		t.Errorf("Assessment = %s, want %s", decoded.Assessment, original.Assessment)
	}
}

func TestTopicFileJSONFieldNames(t *testing.T) {
	topic := TopicFile{
		ID:             "test",
		EstimatedHours: 5,
		ModulePlan:     []ModulePlanEntry{{ID: "m1", Order: 1}},
	}

	data, err := json.Marshal(topic)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	expectedFields := []string{
		"id", "title", "description", "difficulty", "estimated_hours",
		"tags", "prerequisites", "related_topics", "source_urls",
		"generated_at", "version", "module_plan",
	}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("missing expected JSON field %q", field)
		}
	}
}

func TestModuleFileJSONFieldNames(t *testing.T) {
	mod := ModuleFile{
		ID:                 "test",
		LearningObjectives: []string{"obj"},
		Assessment:         json.RawMessage(`{}`),
	}

	data, err := json.Marshal(mod)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	expectedFields := []string{
		"id", "title", "description", "order",
		"learning_objectives", "estimated_minutes", "assessment",
	}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("missing expected JSON field %q", field)
		}
	}
}
