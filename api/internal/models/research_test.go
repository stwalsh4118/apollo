package models_test

import (
	"encoding/json"
	"testing"

	"github.com/sean/apollo/api/internal/models"
)

func TestResearchJobJSONRoundTrip(t *testing.T) {
	progress := models.ResearchProgress{
		CurrentPass:      2,
		TotalPasses:      4,
		ModulesPlanned:   5,
		ModulesCompleted: 1,
		ConceptsFound:    12,
		PassDescriptions: map[int]string{1: "Survey", 2: "Deep Dive"},
	}

	progressJSON, err := json.Marshal(progress)
	if err != nil {
		t.Fatalf("marshal progress: %v", err)
	}

	job := models.ResearchJob{
		ID:        "test-job-1",
		RootTopic: "Go Concurrency",
		Status:    models.ResearchStatusResearching,
		Progress:  progressJSON,
		StartedAt: "2026-02-19T08:00:00Z",
	}

	data, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("marshal job: %v", err)
	}

	var decoded models.ResearchJob
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal job: %v", err)
	}

	if decoded.ID != job.ID {
		t.Fatalf("expected ID %q, got %q", job.ID, decoded.ID)
	}

	if decoded.RootTopic != job.RootTopic {
		t.Fatalf("expected RootTopic %q, got %q", job.RootTopic, decoded.RootTopic)
	}

	if decoded.StartedAt != job.StartedAt {
		t.Fatalf("expected StartedAt %q, got %q", job.StartedAt, decoded.StartedAt)
	}

	if decoded.Status != models.ResearchStatusResearching {
		t.Fatalf("expected status %q, got %q", models.ResearchStatusResearching, decoded.Status)
	}

	var decodedProgress models.ResearchProgress
	if err := json.Unmarshal(decoded.Progress, &decodedProgress); err != nil {
		t.Fatalf("unmarshal progress: %v", err)
	}

	if decodedProgress.CurrentPass != 2 {
		t.Fatalf("expected current_pass 2, got %d", decodedProgress.CurrentPass)
	}

	if decodedProgress.ConceptsFound != 12 {
		t.Fatalf("expected concepts_found 12, got %d", decodedProgress.ConceptsFound)
	}
}

func TestResearchProgressOmitsEmptyFields(t *testing.T) {
	job := models.ResearchJob{
		ID:        "test-job-2",
		RootTopic: "Rust Basics",
		Status:    models.ResearchStatusQueued,
	}

	data, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("marshal job: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	if _, ok := raw["progress"]; ok {
		t.Fatal("expected progress to be omitted when nil")
	}

	if _, ok := raw["error"]; ok {
		t.Fatal("expected error to be omitted when empty")
	}

	if _, ok := raw["started_at"]; ok {
		t.Fatal("expected started_at to be omitted when empty")
	}
}

func TestResearchStatusConstants(t *testing.T) {
	// Verify status constants match the CHECK constraint values in the migration.
	expected := []string{"queued", "researching", "resolving", "published", "failed", "cancelled"}
	actual := []string{
		models.ResearchStatusQueued,
		models.ResearchStatusResearching,
		models.ResearchStatusResolving,
		models.ResearchStatusPublished,
		models.ResearchStatusFailed,
		models.ResearchStatusCancelled,
	}

	for i, exp := range expected {
		if actual[i] != exp {
			t.Fatalf("status constant %d: expected %q, got %q", i, exp, actual[i])
		}
	}
}

func TestIsTerminalStatus(t *testing.T) {
	terminal := []string{
		models.ResearchStatusPublished,
		models.ResearchStatusFailed,
		models.ResearchStatusCancelled,
	}
	for _, s := range terminal {
		if !models.IsTerminalStatus(s) {
			t.Fatalf("expected %q to be terminal", s)
		}
	}

	nonTerminal := []string{
		models.ResearchStatusQueued,
		models.ResearchStatusResearching,
		models.ResearchStatusResolving,
	}
	for _, s := range nonTerminal {
		if models.IsTerminalStatus(s) {
			t.Fatalf("expected %q to be non-terminal", s)
		}
	}
}

func TestCLIResponseJSONParsing(t *testing.T) {
	raw := `{
		"type": "result",
		"session_id": "abc-123",
		"result": "some text output",
		"is_error": false
	}`

	var resp models.CLIResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal CLI response: %v", err)
	}

	if resp.SessionID != "abc-123" {
		t.Fatalf("expected session_id %q, got %q", "abc-123", resp.SessionID)
	}

	if resp.IsError {
		t.Fatal("expected is_error false")
	}

	if resp.Result != "some text output" {
		t.Fatalf("expected result %q, got %q", "some text output", resp.Result)
	}
}

func TestCLIResponseWithStructuredOutput(t *testing.T) {
	raw := `{
		"type": "result",
		"session_id": "def-456",
		"structured_output": {"topic": "Go", "modules": []},
		"is_error": false
	}`

	var resp models.CLIResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal CLI response: %v", err)
	}

	if resp.StructuredOutput == nil {
		t.Fatal("expected structured_output to be non-nil")
	}

	var parsed map[string]any
	if err := json.Unmarshal(resp.StructuredOutput, &parsed); err != nil {
		t.Fatalf("unmarshal structured_output: %v", err)
	}

	if parsed["topic"] != "Go" {
		t.Fatalf("expected topic %q, got %v", "Go", parsed["topic"])
	}

	if parsed["modules"] == nil {
		t.Fatal("expected modules to be present in structured_output")
	}
}
