# PBI-2: Curriculum CRUD API & Search

[View in Backlog](../backlog.md)

## Overview

Build the REST API layer that serves curriculum data — topics, modules, lessons, concepts, concept references, prerequisites, and relations. Also includes the full-text search endpoint. This is the data access layer that the frontend (PBI 3) and research pipeline (PBI 6) both depend on.

## Problem Statement

The database schema exists (PBI 1) but there's no way to read or write curriculum data over HTTP. The frontend needs endpoints to browse topics, drill into modules and lessons, look up concepts, and search. The research pipeline needs endpoints to store newly generated curricula.

## User Stories

- As a learner, I want to browse all topics in my knowledge pool via an API so that the frontend can display them
- As a learner, I want to drill into a topic's modules, then a module's lessons, then a lesson's full content so that I can read course material
- As a learner, I want to look up any concept and see everywhere it's referenced so that I understand cross-topic connections
- As a learner, I want to search across all my content by keyword so that I can find anything I've studied

## Technical Approach

- Go HTTP server (stdlib `net/http` or a lightweight router like `chi`)
- RESTful JSON API matching PRD section 9.2 (Curriculum) and 9.5 (Search)
- Repository pattern: database queries isolated in a `repository` or `store` package
- Endpoints:
  - `GET /api/topics` — list all topics
  - `GET /api/topics/:id` — topic with modules (no lesson content)
  - `GET /api/topics/:id/full` — topic with full tree
  - `GET /api/modules/:id` — module with lessons
  - `GET /api/lessons/:id` — lesson with full content
  - `GET /api/concepts` — list concepts (paginated, filterable)
  - `GET /api/concepts/:id` — concept with references
  - `GET /api/concepts/:id/references` — lessons teaching/referencing this concept
  - `GET /api/search?q=...` — FTS5 full-text search
- Knowledge graph endpoints from PRD section 9.6:
  - `GET /api/graph` — full topic + concept graph (nodes and edges)
  - `GET /api/graph/topic/:id` — subgraph for a single topic
- Internal write endpoints (not user-facing, used by research pipeline):
  - Store/update topics, modules, lessons, concepts, references, prerequisites

## UX/UI Considerations

N/A — backend API PBI. JSON response shapes should be frontend-friendly (nested objects where appropriate, not just flat database rows).

## Acceptance Criteria

1. All GET endpoints from PRD sections 9.2 and 9.6 return correct JSON responses
2. Topic list returns topics ordered by title with difficulty badge and module count
3. Full topic tree endpoint returns nested modules → lessons → concepts
4. Concept references endpoint returns all lessons that teach or reference a concept
5. FTS5 search returns results across topics, lessons, and concepts with relevance ranking
6. Graph endpoint returns nodes (topics + concepts) and edges (prerequisites, references, relations)
7. Write endpoints store curriculum data correctly with foreign key integrity
8. Pagination works on list endpoints (concepts, search results)

## Dependencies

- **Depends on**: PBI 1 (database schema must exist)
- **External**: Go HTTP router (chi or stdlib)

## Open Questions

- None

## Related Tasks

[View Tasks](./tasks.md)
