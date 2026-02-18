package respond_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sean/apollo/api/internal/respond"
)

func TestJSONSetsStatusCode(t *testing.T) {
	rec := httptest.NewRecorder()
	respond.JSON(rec, http.StatusCreated, map[string]string{"id": "1"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
}

func TestJSONSetsContentType(t *testing.T) {
	rec := httptest.NewRecorder()
	respond.JSON(rec, http.StatusOK, map[string]string{"key": "value"})

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type 'application/json', got %q", ct)
	}
}

func TestJSONEncodesBody(t *testing.T) {
	rec := httptest.NewRecorder()
	respond.JSON(rec, http.StatusOK, map[string]string{"key": "value"})

	var m map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&m); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if m["key"] != "value" {
		t.Fatalf("expected key='value', got %q", m["key"])
	}
}

func TestErrorReturnsErrorShape(t *testing.T) {
	rec := httptest.NewRecorder()
	respond.Error(rec, http.StatusNotFound, "not found")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	var m map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&m); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if m["error"] != "not found" {
		t.Fatalf("expected error='not found', got %q", m["error"])
	}
}
