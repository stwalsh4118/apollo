package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// GraphHandler serves knowledge graph endpoints.
type GraphHandler struct {
	repo repository.GraphRepository
}

// NewGraphHandler creates a GraphHandler.
func NewGraphHandler(repo repository.GraphRepository) *GraphHandler {
	return &GraphHandler{repo: repo}
}

// RegisterRoutes mounts graph routes on the given router.
func (h *GraphHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/graph", func(r chi.Router) {
		r.Get("/", h.getFullGraph)
		r.Get("/topic/{id}", h.getTopicGraph)
	})
}

func (h *GraphHandler) getFullGraph(w http.ResponseWriter, r *http.Request) {
	graph, err := h.repo.GetFullGraph(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get graph")

		return
	}

	respond.JSON(w, http.StatusOK, graph)
}

func (h *GraphHandler) getTopicGraph(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	graph, err := h.repo.GetTopicGraph(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get topic graph")

		return
	}

	if graph == nil {
		respond.Error(w, http.StatusNotFound, "topic not found")

		return
	}

	respond.JSON(w, http.StatusOK, graph)
}
