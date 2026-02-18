package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// ConceptHandler serves concept endpoints.
type ConceptHandler struct {
	repo repository.ConceptRepository
}

// NewConceptHandler creates a ConceptHandler.
func NewConceptHandler(repo repository.ConceptRepository) *ConceptHandler {
	return &ConceptHandler{repo: repo}
}

// RegisterRoutes mounts concept routes on the given router.
func (h *ConceptHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/concepts", func(r chi.Router) {
		r.Get("/", h.listConcepts)
		r.Get("/{id}", h.getConceptByID)
		r.Get("/{id}/references", h.getConceptReferences)
	})
}

func (h *ConceptHandler) listConcepts(w http.ResponseWriter, r *http.Request) {
	params := models.ParsePagination(r)
	topicID := r.URL.Query().Get("topic")

	result, err := h.repo.ListConcepts(r.Context(), params, topicID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to list concepts")

		return
	}

	respond.JSON(w, http.StatusOK, result)
}

func (h *ConceptHandler) getConceptByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	concept, err := h.repo.GetConceptByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get concept")

		return
	}

	if concept == nil {
		respond.Error(w, http.StatusNotFound, "concept not found")

		return
	}

	respond.JSON(w, http.StatusOK, concept)
}

func (h *ConceptHandler) getConceptReferences(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	refs, err := h.repo.GetConceptReferences(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get concept references")

		return
	}

	if refs == nil {
		respond.Error(w, http.StatusNotFound, "concept not found")

		return
	}

	respond.JSON(w, http.StatusOK, refs)
}
