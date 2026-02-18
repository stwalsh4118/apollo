# Database API Specification

## Package

`github.com/sean/apollo/api/internal/database`

## Connection Management

```go
// Handle wraps a database connection for Apollo services.
type Handle struct {
    DB *sql.DB
}

// Open creates the SQLite connection, runs migrations, and performs a health check.
func Open(ctx context.Context, databasePath string, logger zerolog.Logger) (*Handle, error)

// Close releases database resources. Safe to call on nil Handle.
func (h *Handle) Close() error

// HealthCheck validates the database connection.
func HealthCheck(ctx context.Context, db *sql.DB) error
```

## Configuration

| Constant | Value | Notes |
|----------|-------|-------|
| Driver | `sqlite` | `modernc.org/sqlite` (pure Go, no CGO) |
| DSN Pragmas | `?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)` | Applied per-connection via DSN |
| Journal Mode | WAL | Persistent pragma, set once via ExecContext |
| Max Open Conns | 1 | SQLite write serialization |
| Max Idle Conns | 1 | Matches open conns |
| Dir Permissions | 0750 | For `MkdirAll` on database directory |

## Migration System

- **Location**: `api/migrations/*.sql` embedded via `embed.FS`
- **Package**: `github.com/sean/apollo/api/migrations` exports `Files embed.FS`
- **Tracking table**: `schema_migrations(id TEXT PK, applied_at TEXT)`
- **Idempotency**: Each migration checked against tracking table before execution
- **Transactions**: Each migration runs in a single transaction with its tracking record

## Schema Tables

| Table | Primary Key | Foreign Keys |
|-------|-------------|--------------|
| `topics` | `id TEXT` | `parent_topic_id -> topics(id)` |
| `modules` | `id TEXT` | `topic_id -> topics(id) ON DELETE CASCADE` |
| `lessons` | `id TEXT` | `module_id -> modules(id) ON DELETE CASCADE` |
| `concepts` | `id TEXT` | `defined_in_lesson -> lessons(id)`, `defined_in_topic -> topics(id)` |
| `concept_references` | `(concept_id, lesson_id)` | Both cascade on delete |
| `topic_prerequisites` | `(topic_id, prerequisite_topic_id)` | Both cascade, self-ref check |
| `topic_relations` | `(topic_a, topic_b)` | Both cascade, self-ref check |
| `expansion_queue` | `id INTEGER AUTOINCREMENT` | `requested_by_topic -> topics(id)` |
| `research_jobs` | `id TEXT` | None |
| `learning_progress` | `lesson_id TEXT` | `lesson_id -> lessons(id) ON DELETE CASCADE` |
| `concept_retention` | `concept_id TEXT` | `concept_id -> concepts(id) ON DELETE CASCADE` |
| `search_index` | FTS5 virtual table | `entity_type`, `entity_id UNINDEXED`, `title`, `body` |

## JSON Columns

All JSON columns use `TEXT` type with `json_valid()` CHECK constraints. Query with `json_extract()`.

Key JSON columns: `topics.tags`, `topics.source_urls`, `modules.learning_objectives`, `modules.assessment`, `lessons.content`, `lessons.examples`, `lessons.exercises`, `lessons.review_questions`, `concepts.aliases`, `research_jobs.progress`.

## Indexes

```sql
idx_modules_topic_id, idx_lessons_module_id, idx_concepts_defined_in_topic,
idx_concept_references_lesson_id, idx_topic_prerequisites_prereq,
idx_topic_relations_topic_b, idx_expansion_queue_status,
idx_expansion_queue_topic_id, idx_research_jobs_status,
idx_learning_progress_status, idx_concept_retention_next_review
```
