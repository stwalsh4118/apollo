package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/respond"
)

// ModuleHandler serves module endpoints.
type ModuleHandler struct {
	repo repository.ModuleRepository
}

// NewModuleHandler creates a ModuleHandler.
func NewModuleHandler(repo repository.ModuleRepository) *ModuleHandler {
	return &ModuleHandler{repo: repo}
}

// RegisterRoutes mounts module routes on the given router.
func (h *ModuleHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/modules/{id}", h.getModuleByID)
}

func (h *ModuleHandler) getModuleByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	module, err := h.repo.GetModuleByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get module")

		return
	}

	if module == nil {
		respond.Error(w, http.StatusNotFound, "module not found")

		return
	}

	respond.JSON(w, http.StatusOK, module)
}
