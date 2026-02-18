package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/internal/database"
	"github.com/sean/apollo/api/internal/server"
)

func setupTestServer(t *testing.T) *server.Server {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	logger := zerolog.New(os.Stderr).Level(zerolog.Disabled)

	handle, err := database.Open(context.Background(), dbPath, logger)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() { _ = handle.Close() })

	return server.New(handle, logger)
}

func TestNewReturnsServer(t *testing.T) {
	srv := setupTestServer(t)
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestRouterReturnsHandler(t *testing.T) {
	srv := setupTestServer(t)
	router := srv.Router()

	if router == nil {
		t.Fatal("expected non-nil router")
	}
}

func TestHealthEndpointReturnsOK(t *testing.T) {
	srv := setupTestServer(t)
	router := srv.Router()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("expected status 'ok', got %q", resp["status"])
	}
}

func TestHealthEndpointSetsContentType(t *testing.T) {
	srv := setupTestServer(t)
	router := srv.Router()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type 'application/json', got %q", ct)
	}
}

func TestRecoveryMiddlewareCatchesPanics(t *testing.T) {
	srv := setupTestServer(t)
	router := srv.Router()

	// Mount a panicking handler on the router to exercise the Recoverer.
	router.Get("/api/test-panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test-panic", nil)
	rec := httptest.NewRecorder()

	// Should not propagate panic — Recoverer should return 500.
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestHealthEndpointWithClosedDB(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	logger := zerolog.New(os.Stderr).Level(zerolog.Disabled)

	handle, err := database.Open(context.Background(), dbPath, logger)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	srv := server.New(handle, logger)
	router := srv.Router()

	// Close the database to simulate an unhealthy state.
	_ = handle.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got %q", resp["status"])
	}
}
