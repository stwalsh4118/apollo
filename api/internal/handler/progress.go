package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// ProgressHandler serves learning progress endpoints.
type ProgressHandler struct {
	repo repository.ProgressRepository
}

// NewProgressHandler creates a ProgressHandler.
func NewProgressHandler(repo repository.ProgressRepository) *ProgressHandler {
	return &ProgressHandler{repo: repo}
}

// RegisterRoutes mounts progress routes on the given router.
func (h *ProgressHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/progress/topics/{id}", h.getTopicProgress)
	r.Put("/api/progress/lessons/{id}", h.updateLessonProgress)
	r.Get("/api/progress/summary", h.getProgressSummary)
}

func (h *ProgressHandler) getTopicProgress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tp, err := h.repo.GetTopicProgress(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, "topic not found")

			return
		}

		respond.Error(w, http.StatusInternalServerError, "failed to get topic progress")

		return
	}

	respond.JSON(w, http.StatusOK, tp)
}

func (h *ProgressHandler) updateLessonProgress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input models.UpdateProgressInput
	if !decodeJSON(w, r, &input) {
		return
	}

	switch input.Status {
	case models.ProgressStatusNotStarted, models.ProgressStatusInProgress, models.ProgressStatusCompleted:
		// valid
	default:
		respond.Error(w, http.StatusBadRequest, "status must be not_started, in_progress, or completed")

		return
	}

	lp, err := h.repo.UpdateLessonProgress(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, "lesson not found")

			return
		}

		respond.Error(w, http.StatusInternalServerError, "failed to update lesson progress")

		return
	}

	respond.JSON(w, http.StatusOK, lp)
}

func (h *ProgressHandler) getProgressSummary(w http.ResponseWriter, r *http.Request) {
	ps, err := h.repo.GetProgressSummary(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get progress summary")

		return
	}

	respond.JSON(w, http.StatusOK, ps)
}
