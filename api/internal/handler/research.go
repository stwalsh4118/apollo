package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// CancelFunc cancels a running research job by ID.
type CancelFunc func(jobID string)

// ResearchHandler serves research job endpoints.
type ResearchHandler struct {
	repo     repository.ResearchJobRepository
	cancelFn CancelFunc
}

// NewResearchHandler creates a ResearchHandler.
func NewResearchHandler(repo repository.ResearchJobRepository, cancelFn CancelFunc) *ResearchHandler {
	return &ResearchHandler{repo: repo, cancelFn: cancelFn}
}

// RegisterRoutes mounts research routes on the given router.
func (h *ResearchHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/research", h.createJob)
	r.Get("/api/research/jobs", h.listJobs)
	r.Get("/api/research/jobs/{id}", h.getJob)
	r.Post("/api/research/jobs/{id}/cancel", h.cancelJob)
}

func (h *ResearchHandler) createJob(w http.ResponseWriter, r *http.Request) {
	var input models.CreateResearchJobInput
	if !decodeJSON(w, r, &input) {
		return
	}

	if input.Topic == "" {
		respond.Error(w, http.StatusBadRequest, "topic is required")

		return
	}

	job, err := h.repo.CreateJob(r.Context(), input)
	if err != nil {
		writeError(w, err)

		return
	}

	respond.JSON(w, http.StatusCreated, job)
}

func (h *ResearchHandler) listJobs(w http.ResponseWriter, r *http.Request) {
	params := models.ParsePagination(r)

	list, err := h.repo.ListJobs(r.Context(), params)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to list research jobs")

		return
	}

	respond.JSON(w, http.StatusOK, list)
}

func (h *ResearchHandler) getJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	job, err := h.repo.GetJobByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, "research job not found")

			return
		}

		respond.Error(w, http.StatusInternalServerError, "failed to get research job")

		return
	}

	respond.JSON(w, http.StatusOK, job)
}

func (h *ResearchHandler) cancelJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	job, err := h.repo.GetJobByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, "research job not found")

			return
		}

		respond.Error(w, http.StatusInternalServerError, "failed to get research job")

		return
	}

	if models.IsTerminalStatus(job.Status) {
		respond.Error(w, http.StatusBadRequest, "job is already in a terminal state")

		return
	}

	// Set cancelled status first, then signal the worker to stop.
	// The worker checks job status before writing its own terminal state.
	if err := h.repo.UpdateJobStatus(r.Context(), id, models.ResearchStatusCancelled, ""); err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to cancel research job")

		return
	}

	if h.cancelFn != nil {
		h.cancelFn(id)
	}

	updated, err := h.repo.GetJobByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to retrieve cancelled job")

		return
	}

	respond.JSON(w, http.StatusOK, updated)
}
