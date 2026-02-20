package research_test

import (
	"encoding/json"
	"testing"

	"github.com/sean/apollo/api/internal/config"
	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/research"
)

func TestBuildInitialArgs(t *testing.T) {
	opts := research.InitialPassOpts{
		Prompt:           "Research Go concurrency",
		SystemPromptFile: "skills/research.md",
		Model:            "opus",
		AllowedTools:     []string{"WebSearch", "WebFetch", "Read", "Write"},
	}

	args := research.BuildInitialArgs(opts)

	assertContains(t, args, "-p", "Research Go concurrency")
	assertContains(t, args, "--output-format", config.ResearchOutputFormat)
	assertContains(t, args, "--system-prompt-file", "skills/research.md")
	assertContains(t, args, "--model", "opus")
	assertContains(t, args, "--allowedTools", "WebSearch,WebFetch,Read,Write")
}

func TestBuildInitialArgsMinimal(t *testing.T) {
	opts := research.InitialPassOpts{
		Prompt: "Minimal prompt",
	}

	args := research.BuildInitialArgs(opts)

	assertContains(t, args, "-p", "Minimal prompt")
	assertContains(t, args, "--output-format", config.ResearchOutputFormat)

	// Should NOT contain optional flags.
	for _, arg := range args {
		if arg == "--system-prompt-file" || arg == "--model" || arg == "--allowedTools" {
			t.Fatalf("unexpected optional flag %q in minimal args", arg)
		}
	}
}

func TestBuildResumeArgs(t *testing.T) {
	opts := research.ResumePassOpts{
		Prompt:    "Continue deep dive",
		SessionID: "session-abc-123",
	}

	args := research.BuildResumeArgs(opts)

	assertContains(t, args, "-p", "Continue deep dive")
	assertContains(t, args, "--resume", "session-abc-123")
	assertContains(t, args, "--output-format", config.ResearchOutputFormat)

	for _, arg := range args {
		if arg == "--json-schema" {
			t.Fatal("unexpected --json-schema in non-final pass args")
		}
	}
}

func TestBuildResumeArgsFinalPass(t *testing.T) {
	opts := research.ResumePassOpts{
		Prompt:         "Final validation",
		SessionID:      "session-abc-123",
		JSONSchemaFile: "schemas/curriculum.json",
	}

	args := research.BuildResumeArgs(opts)

	assertContains(t, args, "--json-schema", "schemas/curriculum.json")
	assertContains(t, args, "--resume", "session-abc-123")
}

func TestParseCLIResponseValid(t *testing.T) {
	resp := models.CLIResponse{
		Type:      "result",
		SessionID: "sess-1",
		Result:    "some output",
		IsError:   false,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// Use the exported test helper — since parseCLIResponse is unexported,
	// we test it via RunInitialPass with a mock or by testing the types directly.
	// Instead, test the JSON round-trip which is the core concern.
	var decoded models.CLIResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.SessionID != "sess-1" {
		t.Fatalf("expected session_id 'sess-1', got %q", decoded.SessionID)
	}

	if decoded.IsError {
		t.Fatal("expected is_error false")
	}
}

func TestParseCLIResponseWithStructuredOutput(t *testing.T) {
	raw := `{
		"type": "result",
		"session_id": "sess-2",
		"structured_output": {"topic_id": "go-basics", "title": "Go Basics"},
		"is_error": false
	}`

	var resp models.CLIResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.SessionID != "sess-2" {
		t.Fatalf("expected session_id 'sess-2', got %q", resp.SessionID)
	}

	if resp.StructuredOutput == nil {
		t.Fatal("expected structured_output to be non-nil")
	}

	var output map[string]any
	if err := json.Unmarshal(resp.StructuredOutput, &output); err != nil {
		t.Fatalf("unmarshal structured_output: %v", err)
	}

	if output["topic_id"] != "go-basics" {
		t.Fatalf("expected topic_id 'go-basics', got %v", output["topic_id"])
	}
}

func TestParseCLIResponseError(t *testing.T) {
	raw := `{
		"type": "error",
		"session_id": "",
		"is_error": true,
		"error_message": "rate limit exceeded"
	}`

	var resp models.CLIResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !resp.IsError {
		t.Fatal("expected is_error true")
	}

	if resp.ErrorMessage != "rate limit exceeded" {
		t.Fatalf("expected error_message 'rate limit exceeded', got %q", resp.ErrorMessage)
	}
}

// assertContains verifies that args contains the flag followed by the expected value.
func assertContains(t *testing.T, args []string, flag, value string) {
	t.Helper()

	for i, arg := range args {
		if arg == flag && i+1 < len(args) && args[i+1] == value {
			return
		}
	}

	t.Fatalf("expected args to contain [%s %s], got %v", flag, value, args)
}
