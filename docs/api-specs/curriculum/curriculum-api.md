# Curriculum API Specification

## Packages

- `github.com/sean/apollo/api/internal/handler` — HTTP handlers
- `github.com/sean/apollo/api/internal/repository` — Data access layer
- `github.com/sean/apollo/api/internal/models` — Request/response types
- `github.com/sean/apollo/api/internal/respond` — JSON response helpers
- `github.com/sean/apollo/api/internal/server` — Server wiring

## REST Endpoints

### Topics

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/topics` | `TopicHandler.listTopics` | List all topics ordered by title |
| GET | `/api/topics/{id}` | `TopicHandler.getTopicByID` | Topic detail with modules |
| GET | `/api/topics/{id}/full` | `TopicHandler.getTopicFull` | Full nested tree (modules > lessons > concepts) |
| POST | `/api/topics` | `WriteHandler.createTopic` | Create topic (201) |
| PUT | `/api/topics/{id}` | `WriteHandler.updateTopic` | Update topic (200) |

### Modules

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/modules/{id}` | `ModuleHandler.getModuleByID` | Module detail with lesson list |
| POST | `/api/modules` | `WriteHandler.createModule` | Create module (201) |

### Lessons

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/lessons/{id}` | `LessonHandler.getLessonByID` | Lesson with full content JSON |
| POST | `/api/lessons` | `WriteHandler.createLesson` | Create lesson (201) |

### Concepts

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/concepts` | `ConceptHandler.listConcepts` | Paginated list (?topic= filter, ?page=, ?per_page=) |
| GET | `/api/concepts/{id}` | `ConceptHandler.getConceptByID` | Concept detail with references |
| GET | `/api/concepts/{id}/references` | `ConceptHandler.getConceptReferences` | Lessons referencing this concept |
| POST | `/api/concepts` | `WriteHandler.createConcept` | Create concept (201) |
| POST | `/api/concepts/{id}/references` | `WriteHandler.createConceptReference` | Add concept reference (201) |

### Search

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/search?q=` | `SearchHandler.search` | FTS5 full-text search (paginated) |

### Knowledge Graph

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/graph` | `GraphHandler.getFullGraph` | All nodes and edges |
| GET | `/api/graph/topic/{id}` | `GraphHandler.getTopicGraph` | Topic subgraph |

### Infrastructure

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/health` | `Server.handleHealth` | Database health check |

### Learning Progress

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/progress/topics/{id}` | `ProgressHandler.getTopicProgress` | Per-lesson progress for a topic (200) |
| PUT | `/api/progress/lessons/{id}` | `ProgressHandler.updateLessonProgress` | Update lesson status and notes (200) |
| GET | `/api/progress/summary` | `ProgressHandler.getProgressSummary` | Completion percentage and active topics (200) |

### Relations & Prerequisites

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/api/prerequisites` | `WriteHandler.createPrerequisite` | Add topic prerequisite (201) |
| POST | `/api/relations` | `WriteHandler.createRelation` | Add topic relation (201) |

## Error Responses

| Status | Sentinel | Meaning |
|--------|----------|---------|
| 400 | — | Invalid JSON, missing required fields, invalid FTS5 query |
| 404 | `ErrNotFound` | Entity not found |
| 409 | `ErrDuplicate` | Duplicate primary key |
| 422 | `ErrFKViolation` | Foreign key constraint violation |
| 400 | `ErrCheckViolation` | CHECK constraint violation |
| 500 | — | Internal server error |

## Repository Interfaces

```go
// TopicRepository — api/internal/repository/topic.go
type TopicRepository interface {
    ListTopics(ctx context.Context) ([]models.TopicSummary, error)
    GetTopicByID(ctx context.Context, id string) (*models.TopicDetail, error)
    GetTopicFull(ctx context.Context, id string) (*models.TopicFull, error)
}

// ModuleRepository — api/internal/repository/module.go
type ModuleRepository interface {
    GetModuleByID(ctx context.Context, id string) (*models.ModuleDetail, error)
}

// LessonRepository — api/internal/repository/lesson.go
type LessonRepository interface {
    GetLessonByID(ctx context.Context, id string) (*models.LessonDetail, error)
}

// ConceptRepository — api/internal/repository/concept.go
type ConceptRepository interface {
    ListConcepts(ctx context.Context, params models.PaginationParams, topicID string) (*models.PaginatedResponse[models.ConceptSummary], error)
    GetConceptByID(ctx context.Context, id string) (*models.ConceptDetail, error)
    GetConceptReferences(ctx context.Context, id string) ([]models.ConceptReference, error)
}

// WriteRepository — api/internal/repository/write.go
type WriteRepository interface {
    CreateTopic(ctx context.Context, input models.TopicInput) error
    UpdateTopic(ctx context.Context, id string, input models.TopicInput) error
    CreateModule(ctx context.Context, input models.ModuleInput) error
    CreateLesson(ctx context.Context, input models.LessonInput) error
    CreateConcept(ctx context.Context, input models.ConceptInput) error
    CreateConceptReference(ctx context.Context, conceptID string, input models.ConceptReferenceInput) error
    CreatePrerequisite(ctx context.Context, input models.PrerequisiteInput) error
    CreateRelation(ctx context.Context, input models.RelationInput) error
}

// SearchRepository — api/internal/repository/search.go
type SearchRepository interface {
    Search(ctx context.Context, query string, params models.PaginationParams) (*models.PaginatedResponse[models.SearchResult], error)
}

// ProgressRepository — api/internal/repository/progress.go
type ProgressRepository interface {
    GetTopicProgress(ctx context.Context, topicID string) (*models.TopicProgress, error)
    UpdateLessonProgress(ctx context.Context, lessonID string, input models.UpdateProgressInput) (*models.LessonProgress, error)
    GetProgressSummary(ctx context.Context) (*models.ProgressSummary, error)
}

// GraphRepository — api/internal/repository/graph.go
type GraphRepository interface {
    GetFullGraph(ctx context.Context) (*models.GraphData, error)
    GetTopicGraph(ctx context.Context, topicID string) (*models.GraphData, error)
}
```

## Sentinel Errors

```go
// api/internal/repository/errors.go
var (
    ErrNotFound       = errors.New("not found")
    ErrDuplicate      = errors.New("duplicate entry")
    ErrFKViolation    = errors.New("foreign key violation")
    ErrCheckViolation = errors.New("check constraint violation")
)

// api/internal/repository/search.go
var ErrInvalidQuery = errors.New("invalid search query")
```

## Response Models

```go
// Pagination — api/internal/models/pagination.go
type PaginationParams struct {
    Page    int // default 1
    PerPage int // default 20, max 100
}
func ParsePagination(r *http.Request) PaginationParams

type PaginatedResponse[T any] struct {
    Items   []T `json:"items"`
    Total   int `json:"total"`
    Page    int `json:"page"`
    PerPage int `json:"per_page"`
}

// Topics — api/internal/models/topic.go
type TopicSummary struct { ID, Title, Description, Difficulty, Status string; ModuleCount int; Tags []string }
type TopicDetail struct { /* base fields + Modules []ModuleSummary */ }
type TopicFull struct { /* base fields + Modules []ModuleFull */ }

// Modules — api/internal/models/module.go
type ModuleSummary struct { ID, Title string; SortOrder int; EstimatedMinutes int }
type ModuleDetail struct { /* base + Lessons []LessonSummary */ }
type ModuleFull struct { /* base + Lessons []LessonFull */ }

// Lessons — api/internal/models/lesson.go
type LessonSummary struct { ID, Title string; SortOrder int; EstimatedMinutes int }
type LessonDetail struct { /* all fields + Content, Examples, Exercises, ReviewQuestions json.RawMessage */ }
type LessonFull struct { /* base + Concepts []ConceptSummary */ }

// Concepts — api/internal/models/concept.go
type ConceptSummary struct { ID, Name, DefinedInTopic string; Aliases []string }
type ConceptDetail struct { /* base + References []ConceptReference */ }
type ConceptReference struct { LessonID, LessonTitle, Context string }

// Search — api/internal/models/search.go
type SearchResult struct { EntityType, EntityID, Title, Snippet string }

// Progress — api/internal/models/progress.go
type LessonProgress struct { LessonID, LessonTitle, Status, StartedAt, CompletedAt, Notes string }
type TopicProgress struct { TopicID string; Lessons []LessonProgress }
type ProgressSummary struct { TotalLessons, CompletedLessons int; CompletionPercentage float64; ActiveTopics int }
type UpdateProgressInput struct { Status, Notes string }

// Graph — api/internal/models/graph.go
type GraphNode struct { ID, Label, Type string }
type GraphEdge struct { Source, Target, Type string }
type GraphData struct { Nodes []GraphNode; Edges []GraphEdge }
```

## Input Models

```go
// api/internal/models/input.go
type TopicInput struct {
    ID, Title, Description, Difficulty, Status, ParentTopicID string
    Tags []string; SourceURLs []string; EstimatedHours float64
}
type ModuleInput struct { ID, TopicID, Title, Description string; SortOrder int; EstimatedMinutes int; LearningObjectives []string }
type LessonInput struct { ID, ModuleID, Title, ContentType string; SortOrder int; EstimatedMinutes int; Content, Examples, Exercises, ReviewQuestions json.RawMessage }
type ConceptInput struct { ID, Name, DefinedInTopic, Description, Importance string; Aliases []string }
type ConceptReferenceInput struct { LessonID, Context string }
type PrerequisiteInput struct { TopicID, PrerequisiteTopicID, Priority string }
type RelationInput struct { TopicA, TopicB, RelationType string }
```

## JSON Response Helpers

```go
// api/internal/respond/respond.go
func JSON(w http.ResponseWriter, status int, data any)           // Sets Content-Type, marshals data
func Error(w http.ResponseWriter, status int, message string)    // Returns {"error": "message"}
```

## Server Configuration

- **Router**: `github.com/go-chi/chi/v5`
- **Middleware**: Recoverer, request logger (zerolog), Content-Type: application/json
- **Max request body**: 2 MB (write endpoints)
- **Graceful shutdown**: SIGINT/SIGTERM with context cancellation
- **Port**: Configured via `SERVER_PORT` env var (default: 8080)
