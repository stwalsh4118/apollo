# PBI-11: Knowledge Wiki & Search UI

[View in Backlog](../backlog.md)

## Overview

Build the knowledge wiki frontend — a browsable index of all topics and concepts with concept detail pages, full-text search, and breadcrumb navigation. This is the reference layer of Apollo: when you want to look something up rather than study sequentially, you use the wiki.

## Problem Statement

The course view (PBI 3) is great for sequential learning but not for reference lookups. The user needs to browse all concepts alphabetically, see where a concept is defined and everywhere it's referenced, and search across everything by keyword. The search API exists (PBI 2) but has no UI.

## User Stories

- As a learner, I want to search across my entire knowledge pool by keyword (US-7)
- As a learner, I want to browse an index of all concepts across all curricula (US-6)
- As a learner, I want a concept detail page showing its definition and every reference (US-6)

## Technical Approach

- Routes:
  - `/wiki` — topic index (list of all topics as a browsable list/grid)
  - `/wiki/concepts` — concept index (alphabetical, filterable by topic)
  - `/wiki/concepts/:id` — concept detail page
  - `/search?q=...` — search results page
- Topic index: list all topics with title, difficulty, module count, concept count
- Concept index: alphabetical list with topic membership badges, filterable
- Concept detail page:
  - Name, definition, difficulty
  - Flashcard (front/back)
  - "Defined in" — link to the canonical lesson
  - "Referenced by" — list of all lessons that reference this concept (from `concept_references`)
  - Related concepts (from the same lessons)
  - Conflict indicator if `status: 'conflict'`
- Search results page:
  - Query input with instant results (debounced)
  - Results grouped by type: Topics, Lessons, Concepts
  - Highlight matching terms in snippets
- Breadcrumb navigation: Topic → Module → Lesson (consistent across wiki and course views)

## UX/UI Considerations

- Wiki should feel like a reference tool — clean, scannable, fast
- Concept detail page is the knowledge hub for each atomic concept
- Search should be fast and forgiving (FTS5 handles stemming)
- Breadcrumbs on every page for orientation

## Acceptance Criteria

1. Topic index displays all topics with metadata
2. Concept index lists all concepts alphabetically with topic badges
3. Concept index filterable by topic
4. Concept detail page shows definition, flashcard, canonical lesson link, and all references
5. Conflict status visible on concept detail page
6. Search results page returns results from FTS5 search grouped by type
7. Search has debounced input (no submit button required)
8. Breadcrumb navigation works across wiki pages
9. Wiki navigation accessible from main app navigation

## Dependencies

- **Depends on**: PBI 2 (search API and concept endpoints), PBI 3 (frontend scaffold)
- **External**: None beyond PBI 3's dependencies

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 11`._
