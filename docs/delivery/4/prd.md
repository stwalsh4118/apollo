# PBI-4: Learning Progress & Notes

[View in Backlog](../backlog.md)

## Overview

Add learning progress tracking and personal notes to both the backend API and the course frontend. Users can mark lessons as complete, see progress bars per module, and write notes on any lesson. Concept chips on lessons link to canonical definitions.

## Problem Statement

The course view (PBI 3) renders content but doesn't track whether the user has studied it. There's no way to mark lessons complete, see how far through a module you are, or capture personal insights alongside generated content. The `learning_progress` table exists but has no API or UI.

## User Stories

- As a learner, I want to mark lessons as complete so that I can track what I've studied (US-3)
- As a learner, I want to see progress bars per module so that I know how far I am (US-3)
- As a learner, I want to add personal notes to any lesson so that I capture my own insights (US-10)
- As a learner, I want concept chips on lessons that link to their definitions so that I can explore connections (US-6)

## Technical Approach

- Backend: Implement PRD section 9.3 (Learning Progress) endpoints:
  - `GET /api/progress/topics/:id` — per-lesson progress for a topic
  - `PUT /api/progress/lessons/:id` — update lesson status and notes
  - `GET /api/progress/summary` — dashboard data (completion %, active topics)
- Frontend additions to the course view:
  - "Mark Complete" button per lesson → `PUT /api/progress/lessons/:id`
  - Progress bar per module in the sidebar (completed/total lessons)
  - Personal notes textarea per lesson (auto-save or explicit save)
  - Concept chips: key concepts for the current lesson shown as clickable badges, linking to `/concepts/:id`
  - Completion indicators in module sidebar (per-lesson check marks)

## UX/UI Considerations

- "Mark Complete" button should be prominent but not intrusive — bottom of lesson
- Progress bars in the sidebar give a sense of momentum
- Notes textarea should be inline within the lesson (not a separate page)
- Concept chips appear near the top of the lesson as a row of clickable badges

## Acceptance Criteria

1. `PUT /api/progress/lessons/:id` stores status (`not_started`, `in_progress`, `completed`) and notes
2. `GET /api/progress/topics/:id` returns per-lesson status for all lessons in a topic
3. `GET /api/progress/summary` returns completion percentage and active topic count
4. "Mark Complete" button in lesson view updates status and shows visual confirmation
5. Module sidebar shows completion indicators (checkmark/filled circle per lesson)
6. Progress bar per module shows correct ratio of completed lessons
7. Personal notes persist across page reloads
8. Concept chips render for each concept taught/referenced in the current lesson
9. Clicking a concept chip navigates to the concept detail (even if wiki isn't built yet, the route exists)

## Dependencies

- **Depends on**: PBI 3 (course frontend must exist to add progress UI to)
- **External**: None beyond PBI 3's dependencies

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 4`._
