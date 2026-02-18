package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/internal/database"
	"github.com/sean/apollo/api/internal/handler"
	"github.com/sean/apollo/api/internal/repository"
)

// Server holds dependencies for the HTTP API.
type Server struct {
	db     *database.Handle
	logger zerolog.Logger
}

// New creates a Server with the given dependencies.
func New(db *database.Handle, logger zerolog.Logger) *Server {
	return &Server{
		db:     db,
		logger: logger,
	}
}

// Router builds and returns the configured chi router with all middleware and routes.
func (s *Server) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(s.requestLogger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/api/health", s.handleHealth)

	topicHandler := handler.NewTopicHandler(repository.NewTopicRepository(s.db.DB))
	topicHandler.RegisterRoutes(r)

	moduleHandler := handler.NewModuleHandler(repository.NewModuleRepository(s.db.DB))
	moduleHandler.RegisterRoutes(r)

	lessonHandler := handler.NewLessonHandler(repository.NewLessonRepository(s.db.DB))
	lessonHandler.RegisterRoutes(r)

	conceptHandler := handler.NewConceptHandler(repository.NewConceptRepository(s.db.DB))
	conceptHandler.RegisterRoutes(r)

	writeHandler := handler.NewWriteHandler(repository.NewWriteRepository(s.db.DB))
	writeHandler.RegisterRoutes(r)

	searchHandler := handler.NewSearchHandler(repository.NewSearchRepository(s.db.DB))
	searchHandler.RegisterRoutes(r)

	graphHandler := handler.NewGraphHandler(repository.NewGraphRepository(s.db.DB))
	graphHandler.RegisterRoutes(r)

	return r
}

// requestLogger is zerolog-based request logging middleware.
func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		next.ServeHTTP(ww, r)

		s.logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.Status()).
			Int("bytes", ww.BytesWritten()).
			Dur("duration", time.Since(start)).
			Msg("request completed")
	})
}
