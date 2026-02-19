package schema

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

// --- AC 6: JSON schema validates a correct curriculum object ---

func TestE2E_ValidCompleteCurriculum(t *testing.T) {
	data := readFixture(t, "valid_curriculum.json")
	if err := Validate(data); err != nil {
		t.Fatalf("AC 6: valid curriculum rejected: %v", err)
	}
}

// --- AC 7: JSON schema rejects malformed objects ---

func TestE2E_RejectMissingRequiredID(t *testing.T) {
	data := readFixture(t, "missing_required_id.json")
	err := Validate(data)
	if err == nil {
		t.Fatal("AC 7: expected rejection for missing required id")
	}
	if !strings.Contains(err.Error(), "id") {
		t.Errorf("AC 7: error should mention 'id', got: %v", err)
	}
}

func TestE2E_RejectInvalidDifficulty(t *testing.T) {
	data := readFixture(t, "invalid_difficulty.json")
	err := Validate(data)
	if err == nil {
		t.Fatal("AC 7: expected rejection for invalid difficulty enum")
	}
}

func TestE2E_RejectMissingContentSections(t *testing.T) {
	data := readFixture(t, "missing_content_sections.json")
	err := Validate(data)
	if err == nil {
		t.Fatal("AC 7: expected rejection for empty content sections")
	}
}

func TestE2E_RejectInvalidExerciseType(t *testing.T) {
	data := readFixture(t, "invalid_exercise_type.json")
	err := Validate(data)
	if err == nil {
		t.Fatal("AC 7: expected rejection for invalid exercise type 'hands_on'")
	}
}

// --- AC 8: Go validation function loads the schema and validates/rejects ---

func TestE2E_GoValidateAcceptsValid(t *testing.T) {
	data := readFixture(t, "valid_curriculum.json")
	if err := Validate(data); err != nil {
		t.Fatalf("AC 8: Validate() should accept valid input: %v", err)
	}
}

func TestE2E_GoValidateRejectsInvalid(t *testing.T) {
	data := readFixture(t, "invalid_difficulty.json")
	if err := Validate(data); err == nil {
		t.Fatal("AC 8: Validate() should reject invalid input")
	}
}

func TestE2E_GoValidateErrorHasFieldPath(t *testing.T) {
	data := readFixture(t, "missing_required_id.json")
	err := Validate(data)
	if err == nil {
		t.Fatal("AC 8: expected error")
	}
	// Error should contain a schema validation prefix.
	if !strings.Contains(err.Error(), "schema validation failed") {
		t.Errorf("AC 8: expected 'schema validation failed' prefix, got: %v", err)
	}
}

// --- AC 9: Knowledge pool summary schema ---

func TestE2E_ValidPoolSummary(t *testing.T) {
	data := readFixture(t, "valid_pool_summary.json")
	if err := ValidatePoolSummary(data); err != nil {
		t.Fatalf("AC 9: valid pool summary rejected: %v", err)
	}
}

func TestE2E_ValidPoolSummaryEmpty(t *testing.T) {
	data := readFixture(t, "valid_pool_summary_empty.json")
	if err := ValidatePoolSummary(data); err != nil {
		t.Fatalf("AC 9: empty pool summary rejected: %v", err)
	}
}

func TestE2E_InvalidPoolSummary(t *testing.T) {
	data := readFixture(t, "invalid_pool_summary.json")
	err := ValidatePoolSummary(data)
	if err == nil {
		t.Fatal("AC 9: expected rejection for invalid pool summary (missing existing_concepts and modules field)")
	}
}
