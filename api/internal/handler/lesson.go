package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// LessonHandler serves lesson endpoints.
type LessonHandler struct {
	repo repository.LessonRepository
}

// NewLessonHandler creates a LessonHandler.
func NewLessonHandler(repo repository.LessonRepository) *LessonHandler {
	return &LessonHandler{repo: repo}
}

// RegisterRoutes mounts lesson routes on the given router.
func (h *LessonHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/lessons/{id}", h.getLessonByID)
}

func (h *LessonHandler) getLessonByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	lesson, err := h.repo.GetLessonByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get lesson")

		return
	}

	if lesson == nil {
		respond.Error(w, http.StatusNotFound, "lesson not found")

		return
	}

	respond.JSON(w, http.StatusOK, lesson)
}
