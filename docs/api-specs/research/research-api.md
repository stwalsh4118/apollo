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
