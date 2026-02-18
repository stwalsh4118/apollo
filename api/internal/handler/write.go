package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// WriteHandler serves curriculum write endpoints.
type WriteHandler struct {
	repo repository.WriteRepository
}

// NewWriteHandler creates a WriteHandler.
func NewWriteHandler(repo repository.WriteRepository) *WriteHandler {
	return &WriteHandler{repo: repo}
}

// RegisterRoutes mounts write routes on the given router.
func (h *WriteHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/topics", h.createTopic)
	r.Put("/api/topics/{id}", h.updateTopic)
	r.Post("/api/modules", h.createModule)
	r.Post("/api/lessons", h.createLesson)
	r.Post("/api/concepts", h.createConcept)
	r.Post("/api/concepts/{id}/references", h.createConceptReference)
	r.Post("/api/prerequisites", h.createPrerequisite)
	r.Post("/api/relations", h.createRelation)
}

func (h *WriteHandler) createTopic(w http.ResponseWriter, r *http.Request) {
	var input models.TopicInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.ID == "" || input.Title == "" || input.Status == "" {
		respond.Error(w, http.StatusBadRequest, "id, title, and status are required")

		return
	}

	if err := h.repo.CreateTopic(r.Context(), input); err != nil {
		writeError(w, err)

		return
	}

	respond.JSON(w, http.StatusCreated, input)
}

func (h *WriteHandler) updateTopic(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input models.TopicInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.Title == "" || input.Status == "" {
		respond.Error(w, http.StatusBadRequest, "title and status are required")

		return
	}

	if err := h.repo.UpdateTopic(r.Context(), id, input); err != nil {
		writeError(w, err)

		return
	}

	input.ID = id
	respond.JSON(w, http.StatusOK, input)
}

func (h *WriteHandler) createModule(w http.ResponseWriter, r *http.Request) {
	var input models.ModuleInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.ID == "" || input.TopicID == "" || input.Title == "" {
		respond.Error(w, http.StatusBadRequest, "id, topic_id, and title are required")

		return
	}

	if err := h.repo.CreateModule(r.Context(), input); err != nil {
		writeError(w, err)

		return
	}

	respond.JSON(w, http.StatusCreated, input)
}

func (h *WriteHandler) createLesson(w http.ResponseWriter, r *http.Request) {
	var input models.LessonInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.ID == "" || input.ModuleID == "" || input.Title == "" || len(input.Content) == 0 {
		respond.Error(w, http.StatusBadRequest, "id, module_id, title, and content are required")

		return
	}

	if err := h.repo.CreateLesson(r.Context(), input); err != nil {
		writeError(w, err)

		return
	}

	respond.JSON(w, http.StatusCreated, input)
}

func (h *WriteHandler) createConcept(w http.ResponseWriter, r *http.Request) {
	var input models.ConceptInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.ID == "" || input.Name == "" || input.Definition == "" {
		respond.Error(w, http.StatusBadRequest, "id, name, and definition are required")

		return
	}

	if err := h.repo.CreateConcept(r.Context(), input); err != nil {
		writeError(w, err)

		return
	}

	respond.JSON(w, http.StatusCreated, input)
}

func (h *WriteHandler) createConceptReference(w http.ResponseWriter, r *http.Request) {
	conceptID := chi.URLParam(r, "id")

	var input models.ConceptReferenceInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.LessonID == "" {
		respond.Error(w, http.StatusBadRequest, "lesson_id is required")

		return
	}

	if err := h.repo.CreateConceptReference(r.Context(), conceptID, input); err != nil {
		writeError(w, err)

		return
	}

	respond.JSON(w, http.StatusCreated, input)
}

func (h *WriteHandler) createPrerequisite(w http.ResponseWriter, r *http.Request) {
	var input models.PrerequisiteInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.TopicID == "" || input.PrerequisiteTopicID == "" || input.Priority == "" {
		respond.Error(w, http.StatusBadRequest, "topic_id, prerequisite_topic_id, and priority are required")

		return
	}

	if err := h.repo.CreatePrerequisite(r.Context(), input); err != nil {
		writeError(w, err)

		return
	}

	respond.JSON(w, http.StatusCreated, input)
}

func (h *WriteHandler) createRelation(w http.ResponseWriter, r *http.Request) {
	var input models.RelationInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.TopicA == "" || input.TopicB == "" || input.RelationType == "" {
		respond.Error(w, http.StatusBadRequest, "topic_a, topic_b, and relation_type are required")

		return
	}

	if err := h.repo.CreateRelation(r.Context(), input); err != nil {
		writeError(w, err)

		return
	}

	respond.JSON(w, http.StatusCreated, input)
}

const maxRequestBodySize = 2 * 1024 * 1024 // 2 MB

func decodeJSON(w http.ResponseWriter, r *http.Request, dest any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid JSON: "+err.Error())

		return false
	}

	return true
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrDuplicate) {
		respond.Error(w, http.StatusConflict, err.Error())

		return
	}

	if errors.Is(err, repository.ErrFKViolation) {
		respond.Error(w, http.StatusUnprocessableEntity, err.Error())

		return
	}

	if errors.Is(err, repository.ErrCheckViolation) {
		respond.Error(w, http.StatusBadRequest, err.Error())

		return
	}

	if errors.Is(err, repository.ErrNotFound) {
		respond.Error(w, http.StatusNotFound, err.Error())

		return
	}

	respond.Error(w, http.StatusInternalServerError, "internal server error")
}
