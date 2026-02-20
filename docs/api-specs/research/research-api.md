# Research API Specification

## Package

`github.com/sean/apollo/api/internal/handler` (ResearchHandler)
`github.com/sean/apollo/api/internal/repository` (ResearchJobRepository)

## REST Endpoints

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/api/research` | `ResearchHandler.createJob` | Create research job (201) |
| GET | `/api/research/jobs` | `ResearchHandler.listJobs` | List jobs with pagination (200) |
| GET | `/api/research/jobs/{id}` | `ResearchHandler.getJob` | Get job by ID (200/404) |
| POST | `/api/research/jobs/{id}/cancel` | `ResearchHandler.cancelJob` | Cancel running job (200/400/404) |

## Request/Response Shapes

### POST /api/research

**Request:**
```json
{ "topic": "Go Concurrency", "brief": "optional guidance" }
```

**Response (201):**
```json
{
  "id": "uuid",
  "root_topic": "Go Concurrency",
  "current_topic": "Go Concurrency",
  "status": "queued",
  "started_at": "",
  "completed_at": ""
}
```

### GET /api/research/jobs

**Query params:** `?page=1&per_page=20`

**Response (200):**
```json
{
  "items": [{ "id": "...", "root_topic": "...", "status": "...", "started_at": "...", "completed_at": "..." }],
  "total": 5,
  "page": 1,
  "per_page": 20
}
```

### GET /api/research/jobs/{id}

**Response (200):** Full `ResearchJob` JSON with `progress` field.

### POST /api/research/jobs/{id}/cancel

**Response (200):** Updated `ResearchJob` with `status: "cancelled"`.
**400:** Job already in terminal state.
**404:** Job not found.

## Repository Interface

```go
type ResearchJobRepository interface {
    CreateJob(ctx context.Context, input models.CreateResearchJobInput) (*models.ResearchJob, error)
    GetJobByID(ctx context.Context, id string) (*models.ResearchJob, error)
    ListJobs(ctx context.Context, params models.PaginationParams) (*models.PaginatedResponse[models.ResearchJobSummary], error)
    UpdateJobStatus(ctx context.Context, id string, status string, errorMsg string) error
    UpdateJobProgress(ctx context.Context, id string, progress models.ResearchProgress) error
    UpdateJobCurrentTopic(ctx context.Context, id string, topic string) error
}
```

## Handler Dependencies

```go
type CancelFunc func(jobID string)

func NewResearchHandler(repo ResearchJobRepository, cancelFn CancelFunc) *ResearchHandler
```

The `CancelFunc` is provided by the orchestrator (Task 6-7) during server startup via `Server.SetCancelResearchFunc()`.

## Error Responses

| Status | Condition |
|--------|-----------|
| 400 | Missing/empty topic, cancel on terminal job, invalid JSON |
| 404 | Job not found |
| 500 | Internal server error |

## File-Per-Lesson Pipeline (Internal)

### Directory Structure

Each research job writes to `data/research/<job-id>/`:

```
topic.json                    # Pass 1: metadata, prerequisites, module plan
modules/
  01-<module-slug>/
    module.json               # Module metadata, learning objectives, assessment
    01-<lesson-slug>.json     # Lesson content, concepts, examples
    02-<lesson-slug>.json
  02-<module-slug>/
    module.json
    01-<lesson-slug>.json
```

### File Types (`research/filetypes.go`)

```go
type TopicFile struct { ID, Title, Description, Difficulty string; EstimatedHours float64; Tags, RelatedTopics, SourceURLs []string; Prerequisites PrerequisitesOutput; ModulePlan []ModulePlanEntry; GeneratedAt string; Version int }
type ModulePlanEntry struct { ID, Title, Description string; Order int }
type ModuleFile struct { ID, Title, Description string; Order int; LearningObjectives []string; EstimatedMinutes int; Assessment json.RawMessage }
```

Lesson files use the existing `LessonOutput` struct from `curriculum.go`.

### Constants

```go
const TopicFileName = "topic.json"
const ModulesDirName = "modules"
const ModuleFileBaseName = "module.json"
```

### Assembler (`research/assembler.go`)

```go
func AssembleFromDir(workDir string) (*CurriculumOutput, error)
```

Reads the file tree, sorts by numeric prefix (`01-`, `02-`), assembles into `CurriculumOutput`, validates against the curriculum schema.

### Orchestrator Flow

All 4 passes use `runPass()` (no `runFinalPass`). After Pass 4, the orchestrator calls `AssembleFromDir(workDir)` → marshals to JSON → feeds to `CurriculumIngester.Ingest()`.
