package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/sean/apollo/api/internal/handler"
	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
)

type mockWriteRepo struct {
	returnErr error
}

func (m *mockWriteRepo) CreateTopic(_ context.Context, _ models.TopicInput) error {
	return m.returnErr
}

func (m *mockWriteRepo) UpdateTopic(_ context.Context, _ string, _ models.TopicInput) error {
	return m.returnErr
}

func (m *mockWriteRepo) CreateModule(_ context.Context, _ models.ModuleInput) error {
	return m.returnErr
}

func (m *mockWriteRepo) CreateLesson(_ context.Context, _ models.LessonInput) error {
	return m.returnErr
}

func (m *mockWriteRepo) CreateConcept(_ context.Context, _ models.ConceptInput) error {
	return m.returnErr
}

func (m *mockWriteRepo) CreateConceptReference(_ context.Context, _ string, _ models.ConceptReferenceInput) error {
	return m.returnErr
}

func (m *mockWriteRepo) CreatePrerequisite(_ context.Context, _ models.PrerequisiteInput) error {
	return m.returnErr
}

func (m *mockWriteRepo) CreateRelation(_ context.Context, _ models.RelationInput) error {
	return m.returnErr
}

func TestCreateTopicHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"id":"t1","title":"Go Basics","status":"published"}`
	req := httptest.NewRequest(http.MethodPost, "/api/topics", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateTopicHandlerMissingFields(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"title":"Go Basics"}`
	req := httptest.NewRequest(http.MethodPost, "/api/topics", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateTopicHandlerInvalidJSON(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPost, "/api/topics", strings.NewReader("{invalid"))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateTopicHandlerDuplicate(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{returnErr: repository.ErrDuplicate}).RegisterRoutes(r)

	body := `{"id":"t1","title":"Go Basics","status":"published"}`
	req := httptest.NewRequest(http.MethodPost, "/api/topics", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
}

func TestCreateModuleHandlerFKViolation(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{returnErr: repository.ErrFKViolation}).RegisterRoutes(r)

	body := `{"id":"m1","topic_id":"nonexistent","title":"Intro","sort_order":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/modules", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}

func TestCreateLessonHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"id":"l1","module_id":"m1","title":"Hello","sort_order":1,"content":[{"type":"text"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/lessons", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateConceptHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"id":"c1","name":"Variable","definition":"A storage location"}`
	req := httptest.NewRequest(http.MethodPost, "/api/concepts", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateTopicHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"title":"Updated Title","status":"published"}`
	req := httptest.NewRequest(http.MethodPut, "/api/topics/t1", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateTopicHandlerNotFound(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{returnErr: repository.ErrNotFound}).RegisterRoutes(r)

	body := `{"title":"Updated Title","status":"published"}`
	req := httptest.NewRequest(http.MethodPut, "/api/topics/nonexistent", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestCreateTopicHandlerInternalError(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{returnErr: errors.New("db error")}).RegisterRoutes(r)

	body := `{"id":"t1","title":"Go Basics","status":"published"}`
	req := httptest.NewRequest(http.MethodPost, "/api/topics", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestCreatePrerequisiteHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"topic_id":"t2","prerequisite_topic_id":"t1","priority":"essential"}`
	req := httptest.NewRequest(http.MethodPost, "/api/prerequisites", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateRelationHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"topic_a":"t1","topic_b":"t2","relation_type":"builds_on"}`
	req := httptest.NewRequest(http.MethodPost, "/api/relations", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateConceptReferenceHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"lesson_id":"l1","context":"Used in examples"}`
	req := httptest.NewRequest(http.MethodPost, "/api/concepts/c1/references", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateConceptReferenceHandlerMissingFields(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{}).RegisterRoutes(r)

	body := `{"context":"some context"}`
	req := httptest.NewRequest(http.MethodPost, "/api/concepts/c1/references", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCheckViolationHandler(t *testing.T) {
	r := chi.NewRouter()
	handler.NewWriteHandler(&mockWriteRepo{returnErr: repository.ErrCheckViolation}).RegisterRoutes(r)

	body := `{"id":"t1","title":"Test","status":"invalid_status"}`
	req := httptest.NewRequest(http.MethodPost, "/api/topics", strings.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
