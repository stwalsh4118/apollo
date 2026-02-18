# PBI-10: Spaced Repetition System

[View in Backlog](../backlog.md)

## Overview

Implement the SM-2 spaced repetition system — both the backend algorithm and the flashcard review UI. Concepts automatically enter the review queue when a user completes the lesson that teaches them. The review session presents flashcards, the user self-rates, and the system schedules the next review per the SM-2 algorithm.

## Problem Statement

Apollo generates curricula and the user can study them (PBIs 3-4), but without spaced repetition, knowledge decays. The PRD specifies SM-2 (section 11) with automatic flashcard generation from concept definitions. The `concept_retention` table exists but has no logic or UI.

## User Stories

- As a learner, I want flashcards automatically generated from concepts I've studied so that I have review material without manual effort (US-5)
- As a learner, I want to review flashcards with spaced repetition so that I actually retain what I learn (US-5)

## Technical Approach

- Backend:
  - SM-2 algorithm implementation per PRD section 11.1 (exact pseudocode provided)
  - Concept lifecycle: `new` → `learning` → `reviewing` → `mastered` (PRD 11.2)
  - Trigger: when a lesson is marked complete (PBI 4), all `concepts_taught` by that lesson move from `new` to `learning`, `next_review` set to tomorrow
  - Review API (PRD section 9.4):
    - `GET /api/review/due` — concepts due for review today
    - `GET /api/review/stats` — due today, upcoming, total mastered
    - `POST /api/review/:conceptId` — submit rating (`forgot`, `hard`, `good`, `easy`)
  - Mastery threshold: interval > 90 days → `mastered` status (configurable `MASTERY_THRESHOLD_DAYS`)
  - `forgot` rating resets to `learning`
- Frontend:
  - Review session page (PRD section 11.3):
    - Show count of concepts due today
    - Flashcard: front (question) only → "Show Answer" button → back (answer) revealed
    - Rating buttons: Forgot / Hard / Good / Easy
    - Progress indicator: "3 of 12"
    - After last card: summary (reviewed count, next due dates)
    - "Re-study" link on forgot → navigate to the lesson where the concept is taught
  - Review queue widget (for later use on Dashboard, PBI 13): "X concepts due today" badge

## UX/UI Considerations

- Flashcard UI should be clean and focused — no distractions during review
- Card flip animation for reveal
- Rating buttons color-coded: Forgot (red), Hard (orange), Good (green), Easy (blue)
- Summary screen should feel rewarding — show streak or progress
- "Re-study" link is important for forgot cards — reconnects to the learning material

## Acceptance Criteria

1. SM-2 algorithm produces correct intervals for all rating combinations (unit tested against PRD pseudocode)
2. Completing a lesson triggers `new` → `learning` transition for all taught concepts
3. `GET /api/review/due` returns concepts where `next_review <= today`
4. `POST /api/review/:conceptId` updates interval, ease_factor, review_count, next_review per SM-2
5. `forgot` rating resets concept to `learning` with interval = 1 day
6. Mastery: concepts with interval > 90 days marked as `mastered`
7. Flashcard UI displays front, reveals back on click, accepts rating
8. Progress indicator shows current position in review session
9. Summary screen shows reviewed count and next due dates
10. "Re-study" link navigates to the correct lesson

## Dependencies

- **Depends on**: PBI 4 (lesson completion triggers concept lifecycle), PBI 3 (frontend scaffold)
- **External**: None

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 10`._
