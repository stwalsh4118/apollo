# PBI-13: Dashboard

[View in Backlog](../backlog.md)

## Overview

Build the dashboard — the home screen of Apollo that ties everything together. Topic cards with completion percentages, a "Currently Studying" section, the review queue widget, active research jobs, available expansions, and a concept map thumbnail.

## Problem Statement

Individual features work (course view, research, review, wiki, map) but there's no unified home screen. The user needs a single place to see: what they're studying, what's due for review, what research is running, and what topics are available to expand. PRD section 13.1 defines the Dashboard view.

## User Stories

- As a learner, I want a home screen showing my learning overview so that I know where to pick up
- As a learner, I want to see my review queue count so that I don't forget daily reviews (US-5)
- As a learner, I want to see available prerequisite expansions so that I can grow my knowledge pool on demand (US-8)
- As a learner, I want to see active research progress so that I know what's being generated (US-12)

## Technical Approach

- Dashboard is the `/` route (or `/dashboard`)
- Data sources:
  - `GET /api/progress/summary` — completion %, active topics
  - `GET /api/review/stats` — due today, upcoming, mastered
  - `GET /api/research/jobs` — active research jobs
  - `GET /api/topics` — all topics for cards
  - Expansion queue (new endpoint or filtered from existing)
- Sections per PRD 13.1:
  - **Topic cards**: each topic with title, difficulty badge, completion %, module count
  - **Currently Studying**: topics with in_progress lessons (filtered from progress data)
  - **Review Queue widget**: "X concepts due today" + "Start Review" button
  - **Research Jobs widget**: active jobs with status, or "No active research" state
  - **Available for Expansion**: helpful/deep prerequisites not yet researched, with "Expand" button → `POST /api/research/expand/:topicId`
  - **Concept Map thumbnail**: small preview of the graph, click to expand to full map

## UX/UI Considerations

- Dashboard should feel like a personal learning cockpit — informative but not overwhelming
- Review queue count should be prominent — daily reviews are the retention engine
- Topic cards should be visual and inviting
- Empty states matter: "No topics yet — start your first research!" when the pool is empty
- Responsive layout for different screen sizes

## Acceptance Criteria

1. Dashboard displays topic cards with title, difficulty, completion %, and module count
2. "Currently Studying" section shows topics with in-progress lessons
3. Review queue widget shows due-today count and links to review session
4. Research jobs widget shows active jobs with status (or empty state)
5. Available expansions listed with "Expand" button that triggers research
6. Concept map thumbnail renders and links to full map
7. Empty states display appropriate messaging for new users
8. Dashboard is the landing page of the application

## Dependencies

- **Depends on**: PBI 4 (progress data), PBI 10 (review stats), PBI 7 (research jobs), PBI 3 (frontend scaffold), PBI 12 (concept map for thumbnail)
- **External**: None

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 13`._
