# PBI-7: Research Progress UI

[View in Backlog](../backlog.md)

## Overview

Build the frontend interface for triggering research and monitoring progress. Users can submit a topic for research, watch the agent work in real-time (polling), see job history, and cancel running jobs. This is the visual verification layer for the research pipeline (PBI 6).

## Problem Statement

The research pipeline (PBI 6) runs but the user has no way to trigger it from the UI or see what's happening. Research takes minutes — without a progress view, the user submits a topic and stares at nothing. The PRD specifies a live progress view (section 13.2) with polling.

## User Stories

- As a learner, I want to submit a topic for research from the UI so that I don't need to use curl/API calls (US-1)
- As a learner, I want to see research progress while agents are working so that I know what's happening (US-12)
- As a learner, I want to see a list of all research jobs (active and completed) so that I have history

## Technical Approach

- "New Research" form: topic name + optional brief → `POST /api/research`
- Research jobs list page: polls `GET /api/research/jobs` periodically
- Research job detail view: polls `GET /api/research/jobs/:id` every 5 seconds (per PRD 13.2)
  - Current topic being researched
  - Current pass (survey / deep dive / exercises / validation)
  - Modules planned and completed (progress bar)
  - Prerequisites discovered (grouped by priority)
  - Concepts identified count
  - Elapsed time
- Cancel button → `POST /api/research/jobs/:id/cancel`
- Job status badges: queued (gray), researching (blue/animated), resolving (yellow), published (green), failed (red), cancelled (gray)
- Link to browse the generated curriculum once status is `published`

## UX/UI Considerations

- Progress view should feel live — animated status indicator during research
- Show enough detail to be informative without being overwhelming
- Published jobs link directly to the course view for the generated topic
- Failed jobs show the error message clearly
- "New Research" form should be accessible from the topic list page and/or a navigation item

## Acceptance Criteria

1. "New Research" form submits a topic and shows the job in the progress view
2. Progress view polls every 5 seconds and updates without full page reload
3. Current pass, module progress, prerequisites, and concept count displayed
4. Elapsed time updates in real-time
5. Cancel button cancels a running job and UI reflects the status change
6. Completed jobs link to the course view for the generated topic
7. Failed jobs display the error message
8. Job list shows all jobs with status badges and timestamps

## Dependencies

- **Depends on**: PBI 6 (research pipeline must be running and returning progress), PBI 3 (frontend scaffold)
- **External**: None beyond PBI 3's dependencies

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 7`._
