# PBI-14: Deployment & Configuration

[View in Backlog](../backlog.md)

## Overview

Package Apollo for self-hosted deployment — Dockerfile, Docker Compose, production build pipeline, environment configuration, and structured logging. The end state: `docker compose up` runs the complete application.

## Problem Statement

Apollo works in development but has no production deployment story. The PRD specifies a single Docker container (section 8.4) that serves both the Go API and the React frontend. Configuration needs to be environment-driven. Logging needs to be structured for troubleshooting. The SQLite database needs a backup story.

## User Stories

- As a developer, I want `docker compose up` to run the complete Apollo stack so that deployment is trivial
- As a developer, I want environment-based configuration so that I can customize paths, ports, and limits without code changes

## Technical Approach

- **Dockerfile** (multi-stage):
  - Stage 1: Build React frontend (`pnpm build`)
  - Stage 2: Build Go binary (`go build`)
  - Stage 3: Minimal runtime image with binary + static assets
  - Go server serves React static files (embedded or from disk)
- **Docker Compose**:
  - Single service: `apollo`
  - Volume mount for `./data/` (SQLite database + research working dir)
  - Environment variables for all PRD section 15 settings
  - Health check endpoint
- **Production configuration**:
  - All settings from PRD section 15 configurable via env vars
  - Sensible defaults baked in
  - Validate required settings on startup
- **Structured logging**:
  - zerolog throughout (already set up in PBI 1, ensure consistent use)
  - Request logging middleware with method, path, status, duration
  - Research job lifecycle logging
- **SQLite backup**:
  - Document the backup strategy: `sqlite3 apollo.db ".backup backup.db"`
  - Optional: expose a `/api/admin/backup` endpoint that triggers a safe online backup
- **Build pipeline**:
  - GitHub Actions: lint, test, build Docker image
  - Makefile with standard targets: `build`, `test`, `lint`, `docker`

## UX/UI Considerations

N/A — infrastructure/DevOps PBI.

## Acceptance Criteria

1. `docker compose up` starts Apollo and serves the frontend on the configured port
2. React frontend served by the Go binary (no separate web server needed)
3. All PRD section 15 settings configurable via environment variables
4. SQLite database persisted via Docker volume mount
5. Structured JSON logging for all HTTP requests and research job events
6. Health check endpoint returns OK when database is accessible
7. Multi-stage Dockerfile produces a minimal image
8. SQLite backup procedure documented and tested

## Dependencies

- **Depends on**: All other PBIs (this packages the complete application)
- **External**: Docker, Docker Compose

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 14`._
