package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// TopicHandler serves topic endpoints.
type TopicHandler struct {
	repo repository.TopicRepository
}

// NewTopicHandler creates a TopicHandler.
func NewTopicHandler(repo repository.TopicRepository) *TopicHandler {
	return &TopicHandler{repo: repo}
}

// RegisterRoutes mounts topic routes on the given router.
func (h *TopicHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/topics", func(r chi.Router) {
		r.Get("/", h.listTopics)
		r.Get("/{id}", h.getTopicByID)
		r.Get("/{id}/full", h.getTopicFull)
	})
}

func (h *TopicHandler) listTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := h.repo.ListTopics(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to list topics")

		return
	}

	respond.JSON(w, http.StatusOK, topics)
}

func (h *TopicHandler) getTopicByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	topic, err := h.repo.GetTopicByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get topic")

		return
	}

	if topic == nil {
		respond.Error(w, http.StatusNotFound, "topic not found")

		return
	}

	respond.JSON(w, http.StatusOK, topic)
}

func (h *TopicHandler) getTopicFull(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	topic, err := h.repo.GetTopicFull(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get topic")

		return
	}

	if topic == nil {
		respond.Error(w, http.StatusNotFound, "topic not found")

		return
	}

	respond.JSON(w, http.StatusOK, topic)
}
