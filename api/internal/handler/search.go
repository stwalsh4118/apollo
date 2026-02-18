package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// SearchHandler serves search endpoints.
type SearchHandler struct {
	repo repository.SearchRepository
}

// NewSearchHandler creates a SearchHandler.
func NewSearchHandler(repo repository.SearchRepository) *SearchHandler {
	return &SearchHandler{repo: repo}
}

// RegisterRoutes mounts search routes on the given router.
func (h *SearchHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/search", h.search)
}

func (h *SearchHandler) search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respond.Error(w, http.StatusBadRequest, "q query parameter is required")

		return
	}

	params := models.ParsePagination(r)

	result, err := h.repo.Search(r.Context(), query, params)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidQuery) {
			respond.Error(w, http.StatusBadRequest, "invalid search query syntax")

			return
		}

		respond.Error(w, http.StatusInternalServerError, "search failed")

		return
	}

	respond.JSON(w, http.StatusOK, result)
}
