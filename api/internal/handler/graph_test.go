package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/handler"
	"github.com/sean/apollo/api/internal/models"
)

type mockGraphRepo struct {
	fullGraph  *models.GraphData
	topicGraph *models.GraphData
	returnErr  error
}

func (m *mockGraphRepo) GetFullGraph(_ context.Context) (*models.GraphData, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.fullGraph, nil
}

func (m *mockGraphRepo) GetTopicGraph(_ context.Context, _ string) (*models.GraphData, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	return m.topicGraph, nil
}

func TestGetFullGraphHandler(t *testing.T) {
	mock := &mockGraphRepo{
		fullGraph: &models.GraphData{
			Nodes: []models.GraphNode{
				{ID: "t1", Label: "Go Basics", Type: "topic"},
				{ID: "c1", Label: "Variable", Type: "concept"},
			},
			Edges: []models.GraphEdge{
				{Source: "c1", Target: "t1", Type: "reference"},
			},
		},
	}

	r := chi.NewRouter()
	handler.NewGraphHandler(mock).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/graph", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	nodes, ok := result["nodes"].([]any)
	if !ok {
		t.Fatal("expected nodes array")
	}

	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
}

func TestGetFullGraphErrorHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewGraphHandler(&mockGraphRepo{returnErr: errors.New("db error")}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/graph", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetTopicGraphHandler(t *testing.T) {
	mock := &mockGraphRepo{
		topicGraph: &models.GraphData{
			Nodes: []models.GraphNode{
				{ID: "t1", Label: "Go Basics", Type: "topic"},
			},
			Edges: []models.GraphEdge{},
		},
	}

	r := chi.NewRouter()
	handler.NewGraphHandler(mock).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/graph/topic/t1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestGetTopicGraphNotFoundHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewGraphHandler(&mockGraphRepo{topicGraph: nil}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/graph/topic/nonexistent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
