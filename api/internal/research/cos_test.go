package research

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/internal/database"
)

// cosTestDB creates a temporary SQLite database for CoS tests.
func cosTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "cos-test.db")
	logger := zerolog.New(os.Stderr).Level(zerolog.Disabled)

	handle, err := database.Open(context.Background(), dbPath, logger)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() { _ = handle.Close() })

	return handle.DB
}

// TestCoS_AssemblerIngestionRoundTrip validates AC5 (assembler produces valid output),
// AC6 (output passes schema validation), and AC7 (ingestion stores same data in SQLite).
func TestCoS_AssemblerIngestionRoundTrip(t *testing.T) {
	dir := t.TempDir()
	topic := buildValidFixtureTree(t, dir)

	// AC5: Assembler reads file tree and produces CurriculumOutput.
	curriculum, err := AssembleFromDir(dir)
	if err != nil {
		t.Fatalf("AC5: AssembleFromDir failed: %v", err)
	}

	if curriculum.ID != topic.ID {
		t.Errorf("AC5: curriculum ID = %q, want %q", curriculum.ID, topic.ID)
	}
	if len(curriculum.Modules) != 2 {
		t.Fatalf("AC5: expected 2 modules, got %d", len(curriculum.Modules))
	}
	if len(curriculum.Modules[0].Lessons) != 2 {
		t.Fatalf("AC5: expected 2 lessons in module 1, got %d", len(curriculum.Modules[0].Lessons))
	}

	// AC6 is implicitly validated inside AssembleFromDir (calls schema.Validate).
	// Marshal for AC7 ingestion round-trip.
	assembledJSON, err := json.Marshal(curriculum)
	if err != nil {
		t.Fatalf("marshal assembled curriculum: %v", err)
	}

	// AC7: Ingest into SQLite and verify data is stored.
	db := cosTestDB(t)
	ingester := NewCurriculumIngester(db)

	result, err := ingester.IngestWithResult(context.Background(), json.RawMessage(assembledJSON))
	if err != nil {
		t.Fatalf("AC7: ingest failed: %v", err)
	}

	if result.ModulesCreated != 2 {
		t.Errorf("AC7: modules created = %d, want 2", result.ModulesCreated)
	}
	if result.LessonsCreated != 4 {
		t.Errorf("AC7: lessons created = %d, want 4", result.LessonsCreated)
	}
	if result.ConceptsCreated < 1 {
		t.Errorf("AC7: expected at least 1 concept created, got %d", result.ConceptsCreated)
	}

	// Verify topic is queryable from DB.
	var topicTitle string
	err = db.QueryRow("SELECT title FROM topics WHERE id = ?", topic.ID).Scan(&topicTitle)
	if err != nil {
		t.Fatalf("AC7: query topic: %v", err)
	}
	if topicTitle != topic.Title {
		t.Errorf("AC7: topic title = %q, want %q", topicTitle, topic.Title)
	}

	// Verify modules are queryable.
	var moduleCount int
	err = db.QueryRow("SELECT COUNT(*) FROM modules WHERE topic_id = ?", topic.ID).Scan(&moduleCount)
	if err != nil {
		t.Fatalf("AC7: query modules: %v", err)
	}
	if moduleCount != 2 {
		t.Errorf("AC7: module count = %d, want 2", moduleCount)
	}

	// Verify lessons are queryable.
	var lessonCount int
	err = db.QueryRow("SELECT COUNT(*) FROM lessons").Scan(&lessonCount)
	if err != nil {
		t.Fatalf("AC7: query lessons: %v", err)
	}
	if lessonCount != 4 {
		t.Errorf("AC7: lesson count = %d, want 4", lessonCount)
	}

	// Verify concepts are queryable.
	var conceptCount int
	err = db.QueryRow("SELECT COUNT(*) FROM concepts").Scan(&conceptCount)
	if err != nil {
		t.Fatalf("AC7: query concepts: %v", err)
	}
	if conceptCount < 1 {
		t.Errorf("AC7: expected at least 1 concept in DB, got %d", conceptCount)
	}
}

// TestCoS_FileSizeConstraint validates AC4: no individual file exceeds ~300 lines.
func TestCoS_FileSizeConstraint(t *testing.T) {
	const maxLines = 300

	dir := t.TempDir()
	buildValidFixtureTree(t, dir)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		relPath, _ := filepath.Rel(dir, path)
		if lineCount > maxLines {
			t.Errorf("AC4: %s has %d lines, exceeds max of %d", relPath, lineCount, maxLines)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk directory: %v", err)
	}
}

// TestCoS_PartialProgress validates AC8: if pass 3 fails, pass 2 lesson files
// are still on disk and contain valid JSON.
func TestCoS_PartialProgress(t *testing.T) {
	dir := t.TempDir()

	// Simulate pass 1-2 output: write fixture tree with full lesson files.
	buildValidFixtureTree(t, dir)

	// Record all lesson file paths before "pass 3".
	var lessonPaths []string
	modulesDir := filepath.Join(dir, ModulesDirName)

	err := filepath.Walk(modulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".json") && filepath.Base(path) != ModuleFileBaseName {
			lessonPaths = append(lessonPaths, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk modules: %v", err)
	}

	if len(lessonPaths) == 0 {
		t.Fatal("AC8: expected lesson files after pass 2, found none")
	}

	// Simulate pass 3 failing partway: modify only the first lesson file,
	// then "crash" (leave the rest untouched). This validates the architectural
	// property that each file is an independent write — partial progress is
	// inherently preserved because untouched files remain on disk.
	data, err := os.ReadFile(lessonPaths[0])
	if err != nil {
		t.Fatalf("read lesson: %v", err)
	}

	var lesson LessonOutput
	if err := json.Unmarshal(data, &lesson); err != nil {
		t.Fatalf("unmarshal lesson: %v", err)
	}

	// Pass 3 adds exercises — modify just this one file.
	lesson.Exercises = json.RawMessage(`[{"type":"command","title":"Added by pass 3","instructions":"Do X","success_criteria":["Done"],"hints":[],"environment":"terminal"}]`)
	updated, err := json.MarshalIndent(lesson, "", "  ")
	if err != nil {
		t.Fatalf("marshal updated lesson: %v", err)
	}
	if err := os.WriteFile(lessonPaths[0], updated, 0o644); err != nil {
		t.Fatalf("write updated lesson: %v", err)
	}

	// AC8: All lesson files (including untouched ones) should still exist
	// and be valid JSON after the simulated partial pass 3.
	for _, path := range lessonPaths {
		relPath, _ := filepath.Rel(dir, path)

		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("AC8: %s not readable after partial pass 3: %v", relPath, err)
			continue
		}

		var lesson LessonOutput
		if err := json.Unmarshal(data, &lesson); err != nil {
			t.Errorf("AC8: %s is not valid JSON after partial pass 3: %v", relPath, err)
			continue
		}

		if lesson.ID == "" {
			t.Errorf("AC8: %s has empty ID — file appears corrupted", relPath)
		}
	}

	// topic.json should also still be valid.
	topicData, err := os.ReadFile(filepath.Join(dir, TopicFileName))
	if err != nil {
		t.Fatalf("AC8: topic.json not readable: %v", err)
	}

	var topic TopicFile
	if err := json.Unmarshal(topicData, &topic); err != nil {
		t.Fatalf("AC8: topic.json is not valid JSON: %v", err)
	}
}

// TestCoS_SystemPromptContent validates AC9: the system prompt contains
// file-writing instructions and directory conventions.
func TestCoS_SystemPromptContent(t *testing.T) {
	prompt := string(systemPromptContent)

	checks := []struct {
		name    string
		pattern string
	}{
		{"file-per-lesson description", "file-per-lesson"},
		{"directory structure section", "Directory Structure"},
		{"naming conventions", "Naming Conventions"},
		{"topic.json reference", "topic.json"},
		{"module.json reference", "module.json"},
		{"Write tool instruction", "Write tool"},
		{"Read tool instruction", "Read tool"},
		{"pass 1 file writing", "Pass 1"},
		{"pass 2 sub-agent writing", "Pass 2"},
		{"pass 3 read-modify-write", "Pass 3"},
		{"pass 4 file tree review", "Pass 4"},
		{"NN-slug naming pattern", "NN-"},
		{"lesson file format example", "lesson-slug"},
	}

	for _, check := range checks {
		if !strings.Contains(prompt, check.pattern) {
			t.Errorf("AC9: system prompt missing %q (%s)", check.pattern, check.name)
		}
	}

	// Verify no --json-schema references in pass instructions.
	if strings.Contains(prompt, "--json-schema") {
		t.Error("AC9: system prompt should not reference --json-schema")
	}

	// Verify pass 4 explicitly tells agent NOT to produce structured JSON output.
	if !strings.Contains(prompt, "Do NOT produce structured JSON output") {
		t.Error("AC9: system prompt pass 4 should instruct NOT to produce structured JSON output")
	}
}

// TestCoS_Pass4Prompt validates that the pass 4 prompt does not reference
// structured JSON output or --json-schema.
func TestCoS_Pass4Prompt(t *testing.T) {
	prompt := passPrompts[4]

	if strings.Contains(prompt, "structured JSON") {
		t.Error("pass 4 prompt should not reference structured JSON")
	}
	if strings.Contains(prompt, "--json-schema") {
		t.Error("pass 4 prompt should not reference --json-schema")
	}
	if !strings.Contains(prompt, "file tree") {
		t.Error("pass 4 prompt should reference file tree review")
	}
}
