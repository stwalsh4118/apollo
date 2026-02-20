package research

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFixtureFile writes JSON-encoded data to the given path.
func writeFixtureFile(t *testing.T, path string, v any) {
	t.Helper()

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture %s: %v", path, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for fixture %s: %v", path, err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}

// buildValidFixtureTree creates a complete, schema-valid directory tree in dir.
// Returns the TopicFile used so callers can verify assembled output.
func buildValidFixtureTree(t *testing.T, dir string) TopicFile {
	t.Helper()

	topic := TopicFile{
		ID:             "go-concurrency",
		Title:          "Go Concurrency",
		Description:    "Learn concurrent programming in Go.",
		Difficulty:     "intermediate",
		EstimatedHours: 10,
		Tags:           []string{"go", "concurrency"},
		Prerequisites: PrerequisitesOutput{
			Essential:      []PrerequisiteItem{{TopicID: "go-basics", Reason: "core language"}},
			Helpful:        []PrerequisiteItem{},
			DeepBackground: []PrerequisiteItem{},
		},
		RelatedTopics: []string{"distributed-systems"},
		SourceURLs:    []string{"https://go.dev/doc/effective_go"},
		GeneratedAt:   "2026-02-20T10:00:00Z",
		Version:       1,
		ModulePlan: []ModulePlanEntry{
			{ID: "goroutines", Title: "Goroutines", Description: "Lightweight threads", Order: 1},
			{ID: "channels", Title: "Channels", Description: "Communication primitives", Order: 2},
		},
	}
	writeFixtureFile(t, filepath.Join(dir, TopicFileName), topic)

	// Module 1: goroutines
	mod1Dir := filepath.Join(dir, ModulesDirName, "01-goroutines")
	mod1 := ModuleFile{
		ID:                 "goroutines",
		Title:              "Goroutines",
		Description:        "Lightweight threads in Go",
		Order:              1,
		LearningObjectives: []string{"Understand goroutine lifecycle"},
		EstimatedMinutes:   60,
		Assessment:         json.RawMessage(`{"questions":[{"type":"conceptual","question":"What is a goroutine?","answer":"A lightweight thread managed by the Go runtime.","concepts_tested":["goroutine-basics"]}]}`),
	}
	writeFixtureFile(t, filepath.Join(mod1Dir, ModuleFileBaseName), mod1)

	lesson1 := LessonOutput{
		ID:               "goroutine-basics",
		Title:            "Goroutine Basics",
		Order:            1,
		EstimatedMinutes: 30,
		Content:          json.RawMessage(`{"sections":[{"type":"text","body":"Goroutines are lightweight threads."}]}`),
		ConceptsTaught: []ConceptTaughtOut{{
			ID: "goroutine", Name: "Goroutine",
			Definition: "A lightweight thread of execution.",
			Flashcard:  FlashcardOut{Front: "What is a goroutine?", Back: "A lightweight thread."},
		}},
		ConceptsReferenced: []ConceptRefOut{},
		Examples:           json.RawMessage(`[]`),
		Exercises:          json.RawMessage(`[{"type":"command","title":"Run a goroutine","instructions":"Write a goroutine.","success_criteria":["Output appears"],"hints":[],"environment":"terminal"}]`),
		ReviewQuestions:    json.RawMessage(`[{"question":"What keyword starts a goroutine?","answer":"go","concepts_tested":["goroutine"]}]`),
	}
	writeFixtureFile(t, filepath.Join(mod1Dir, "01-goroutine-basics.json"), lesson1)

	lesson2 := LessonOutput{
		ID:               "goroutine-patterns",
		Title:            "Goroutine Patterns",
		Order:            2,
		EstimatedMinutes: 30,
		Content:          json.RawMessage(`{"sections":[{"type":"text","body":"Common goroutine patterns."}]}`),
		ConceptsTaught: []ConceptTaughtOut{{
			ID: "fan-out", Name: "Fan-Out",
			Definition: "Starting multiple goroutines to handle input.",
			Flashcard:  FlashcardOut{Front: "What is fan-out?", Back: "Starting multiple goroutines."},
		}},
		ConceptsReferenced: []ConceptRefOut{{ID: "goroutine", DefinedIn: "goroutine-basics"}},
		Examples:           json.RawMessage(`[]`),
		Exercises:          json.RawMessage(`[{"type":"build","title":"Build a fan-out","instructions":"Implement fan-out.","success_criteria":["Multiple goroutines run"],"hints":[],"environment":"terminal"}]`),
		ReviewQuestions:    json.RawMessage(`[{"question":"What is fan-out?","answer":"Starting multiple goroutines.","concepts_tested":["fan-out"]}]`),
	}
	writeFixtureFile(t, filepath.Join(mod1Dir, "02-goroutine-patterns.json"), lesson2)

	// Module 2: channels
	mod2Dir := filepath.Join(dir, ModulesDirName, "02-channels")
	mod2 := ModuleFile{
		ID:                 "channels",
		Title:              "Channels",
		Description:        "Communication between goroutines",
		Order:              2,
		LearningObjectives: []string{"Use channels for goroutine communication"},
		EstimatedMinutes:   60,
		Assessment:         json.RawMessage(`{"questions":[{"type":"practical","question":"How do you send a value on a channel?","answer":"Use ch <- value syntax.","concepts_tested":["channel-send"]}]}`),
	}
	writeFixtureFile(t, filepath.Join(mod2Dir, ModuleFileBaseName), mod2)

	lesson3 := LessonOutput{
		ID:               "channel-basics",
		Title:            "Channel Basics",
		Order:            1,
		EstimatedMinutes: 30,
		Content:          json.RawMessage(`{"sections":[{"type":"text","body":"Channels are typed conduits."}]}`),
		ConceptsTaught: []ConceptTaughtOut{{
			ID: "channel", Name: "Channel",
			Definition: "A typed conduit for goroutine communication.",
			Flashcard:  FlashcardOut{Front: "What is a channel?", Back: "A typed conduit."},
		}},
		ConceptsReferenced: []ConceptRefOut{{ID: "goroutine", DefinedIn: "goroutine-basics"}},
		Examples:           json.RawMessage(`[]`),
		Exercises:          json.RawMessage(`[{"type":"command","title":"Create a channel","instructions":"Use make(chan int).","success_criteria":["Channel created"],"hints":[],"environment":"terminal"}]`),
		ReviewQuestions:    json.RawMessage(`[{"question":"How do you create a channel?","answer":"make(chan Type)","concepts_tested":["channel"]}]`),
	}
	writeFixtureFile(t, filepath.Join(mod2Dir, "01-channel-basics.json"), lesson3)

	lesson4 := LessonOutput{
		ID:               "channel-patterns",
		Title:            "Channel Patterns",
		Order:            2,
		EstimatedMinutes: 30,
		Content:          json.RawMessage(`{"sections":[{"type":"text","body":"Select and pipeline patterns."}]}`),
		ConceptsTaught: []ConceptTaughtOut{{
			ID: "select-statement", Name: "Select Statement",
			Definition: "Multiplexes channel operations.",
			Flashcard:  FlashcardOut{Front: "What does select do?", Back: "Multiplexes channel operations."},
		}},
		ConceptsReferenced: []ConceptRefOut{{ID: "channel", DefinedIn: "channel-basics"}},
		Examples:           json.RawMessage(`[]`),
		Exercises:          json.RawMessage(`[{"type":"build","title":"Build a pipeline","instructions":"Chain channels.","success_criteria":["Pipeline works"],"hints":[],"environment":"terminal"}]`),
		ReviewQuestions:    json.RawMessage(`[{"question":"What is select?","answer":"A control structure for channel operations.","concepts_tested":["select-statement"]}]`),
	}
	writeFixtureFile(t, filepath.Join(mod2Dir, "02-channel-patterns.json"), lesson4)

	return topic
}

func TestAssembleFromDir_ValidTree(t *testing.T) {
	dir := t.TempDir()
	topic := buildValidFixtureTree(t, dir)

	curriculum, err := AssembleFromDir(dir)
	if err != nil {
		t.Fatalf("AssembleFromDir: %v", err)
	}

	if curriculum.ID != topic.ID {
		t.Errorf("ID = %q, want %q", curriculum.ID, topic.ID)
	}
	if curriculum.Title != topic.Title {
		t.Errorf("Title = %q, want %q", curriculum.Title, topic.Title)
	}
	if len(curriculum.Modules) != 2 {
		t.Fatalf("Modules count = %d, want 2", len(curriculum.Modules))
	}
	if curriculum.Modules[0].ID != "goroutines" {
		t.Errorf("Modules[0].ID = %q, want %q", curriculum.Modules[0].ID, "goroutines")
	}
	if curriculum.Modules[1].ID != "channels" {
		t.Errorf("Modules[1].ID = %q, want %q", curriculum.Modules[1].ID, "channels")
	}
	if len(curriculum.Modules[0].Lessons) != 2 {
		t.Errorf("Modules[0].Lessons count = %d, want 2", len(curriculum.Modules[0].Lessons))
	}
	if len(curriculum.Modules[1].Lessons) != 2 {
		t.Errorf("Modules[1].Lessons count = %d, want 2", len(curriculum.Modules[1].Lessons))
	}
}

func TestAssembleFromDir_MissingTopicJSON(t *testing.T) {
	dir := t.TempDir()

	_, err := AssembleFromDir(dir)
	if err == nil {
		t.Fatal("expected error for missing topic.json, got nil")
	}

	if want := "read topic.json"; !strings.Contains(err.Error(), want) {
		t.Errorf("error = %q, want to contain %q", err.Error(), want)
	}
}

func TestAssembleFromDir_EmptyModulesDir(t *testing.T) {
	dir := t.TempDir()

	// Write topic.json but create empty modules/.
	topic := TopicFile{ID: "test", Title: "Test", Description: "Test topic.",
		Difficulty: "foundational", EstimatedHours: 1, Tags: []string{},
		Prerequisites: PrerequisitesOutput{Essential: []PrerequisiteItem{}, Helpful: []PrerequisiteItem{}, DeepBackground: []PrerequisiteItem{}},
		RelatedTopics: []string{}, SourceURLs: []string{"https://example.com"},
		GeneratedAt: "2026-01-01T00:00:00Z", Version: 1, ModulePlan: []ModulePlanEntry{},
	}
	writeFixtureFile(t, filepath.Join(dir, TopicFileName), topic)

	if err := os.MkdirAll(filepath.Join(dir, ModulesDirName), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := AssembleFromDir(dir)
	if err == nil {
		t.Fatal("expected error for empty modules directory, got nil")
	}

	if want := "empty"; !strings.Contains(err.Error(), want) {
		t.Errorf("error = %q, want to contain %q", err.Error(), want)
	}
}

func TestAssembleFromDir_MissingModuleJSON(t *testing.T) {
	dir := t.TempDir()

	topic := TopicFile{ID: "test", Title: "Test", Description: "Test topic.",
		Difficulty: "foundational", EstimatedHours: 1, Tags: []string{},
		Prerequisites: PrerequisitesOutput{Essential: []PrerequisiteItem{}, Helpful: []PrerequisiteItem{}, DeepBackground: []PrerequisiteItem{}},
		RelatedTopics: []string{}, SourceURLs: []string{"https://example.com"},
		GeneratedAt: "2026-01-01T00:00:00Z", Version: 1, ModulePlan: []ModulePlanEntry{},
	}
	writeFixtureFile(t, filepath.Join(dir, TopicFileName), topic)

	// Create module dir without module.json.
	modDir := filepath.Join(dir, ModulesDirName, "01-intro")
	if err := os.MkdirAll(modDir, 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := AssembleFromDir(dir)
	if err == nil {
		t.Fatal("expected error for missing module.json, got nil")
	}

	if want := "module.json"; !strings.Contains(err.Error(), want) {
		t.Errorf("error = %q, want to contain %q", err.Error(), want)
	}
}

func TestAssembleFromDir_MalformedLessonJSON(t *testing.T) {
	dir := t.TempDir()

	topic := TopicFile{ID: "test", Title: "Test", Description: "Test topic.",
		Difficulty: "foundational", EstimatedHours: 1, Tags: []string{},
		Prerequisites: PrerequisitesOutput{Essential: []PrerequisiteItem{}, Helpful: []PrerequisiteItem{}, DeepBackground: []PrerequisiteItem{}},
		RelatedTopics: []string{}, SourceURLs: []string{"https://example.com"},
		GeneratedAt: "2026-01-01T00:00:00Z", Version: 1, ModulePlan: []ModulePlanEntry{},
	}
	writeFixtureFile(t, filepath.Join(dir, TopicFileName), topic)

	modDir := filepath.Join(dir, ModulesDirName, "01-intro")
	mod := ModuleFile{ID: "intro", Title: "Intro", Description: "Introduction",
		Order: 1, LearningObjectives: []string{"Learn"}, EstimatedMinutes: 30,
		Assessment: json.RawMessage(`{"questions":[{"type":"conceptual","question":"Q?","answer":"A","concepts_tested":[]}]}`),
	}
	writeFixtureFile(t, filepath.Join(modDir, ModuleFileBaseName), mod)

	// Write malformed JSON as a lesson file.
	if err := os.WriteFile(filepath.Join(modDir, "01-bad.json"), []byte("{invalid json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := AssembleFromDir(dir)
	if err == nil {
		t.Fatal("expected error for malformed lesson JSON, got nil")
	}

	if want := "01-bad.json"; !strings.Contains(err.Error(), want) {
		t.Errorf("error = %q, want to contain %q", err.Error(), want)
	}
}

func TestAssembleFromDir_ModuleSortOrder(t *testing.T) {
	dir := t.TempDir()
	buildValidFixtureTree(t, dir)

	// Rename module dirs to reverse filesystem order.
	modulesDir := filepath.Join(dir, ModulesDirName)
	if err := os.Rename(filepath.Join(modulesDir, "01-goroutines"), filepath.Join(modulesDir, "10-goroutines")); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(modulesDir, "02-channels"), filepath.Join(modulesDir, "05-channels")); err != nil {
		t.Fatal(err)
	}

	curriculum, err := AssembleFromDir(dir)
	if err != nil {
		t.Fatalf("AssembleFromDir: %v", err)
	}

	// channels (05) should come before goroutines (10).
	if curriculum.Modules[0].ID != "channels" {
		t.Errorf("Modules[0].ID = %q, want %q (sorted by prefix)", curriculum.Modules[0].ID, "channels")
	}
	if curriculum.Modules[1].ID != "goroutines" {
		t.Errorf("Modules[1].ID = %q, want %q (sorted by prefix)", curriculum.Modules[1].ID, "goroutines")
	}
}

func TestAssembleFromDir_LessonSortOrder(t *testing.T) {
	dir := t.TempDir()
	buildValidFixtureTree(t, dir)

	// Rename lesson files in module 1 to reverse numeric order.
	mod1Dir := filepath.Join(dir, ModulesDirName, "01-goroutines")
	if err := os.Rename(filepath.Join(mod1Dir, "01-goroutine-basics.json"), filepath.Join(mod1Dir, "10-goroutine-basics.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(mod1Dir, "02-goroutine-patterns.json"), filepath.Join(mod1Dir, "05-goroutine-patterns.json")); err != nil {
		t.Fatal(err)
	}

	curriculum, err := AssembleFromDir(dir)
	if err != nil {
		t.Fatalf("AssembleFromDir: %v", err)
	}

	// goroutine-patterns (05) should come before goroutine-basics (10).
	if curriculum.Modules[0].Lessons[0].ID != "goroutine-patterns" {
		t.Errorf("Lessons[0].ID = %q, want %q", curriculum.Modules[0].Lessons[0].ID, "goroutine-patterns")
	}
	if curriculum.Modules[0].Lessons[1].ID != "goroutine-basics" {
		t.Errorf("Lessons[1].ID = %q, want %q", curriculum.Modules[0].Lessons[1].ID, "goroutine-basics")
	}
}

func TestNumericPrefix(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"01-introduction", 1},
		{"02-basics.json", 2},
		{"10-advanced", 10},
		{"no-prefix", 0},
		{"", 0},
	}

	for _, tt := range tests {
		got := numericPrefix(tt.name)
		if got != tt.want {
			t.Errorf("numericPrefix(%q) = %d, want %d", tt.name, got, tt.want)
		}
	}
}
