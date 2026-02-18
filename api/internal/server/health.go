package server

import (
	"encoding/json"
	"net/http"

	"github.com/sean/apollo/api/internal/database"
)

type healthResponse struct {
	Status string `json:"status"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := database.HealthCheck(r.Context(), s.db.DB); err != nil {
		s.logger.Error().Err(err).Msg("health check failed")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(healthResponse{Status: "error"})

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}
