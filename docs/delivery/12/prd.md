# PBI-12: Concept Map Visualization

[View in Backlog](../backlog.md)

## Overview

Build the interactive concept map — a D3.js force-directed graph that visualizes topics as large nodes, concepts as small nodes, and their connections as typed edges. Users can zoom, filter, click to navigate, and color-code by various dimensions.

## Problem Statement

The knowledge wiki (PBI 11) lets you look up individual concepts, but there's no way to see the big picture — how topics relate, where concept clusters form, where knowledge gaps exist. The PRD specifies a force-directed graph (section 13.1, Concept Map) that gives a bird's-eye view of the entire knowledge pool.

## User Stories

- As a learner, I want a visual map of all my topics and their connections so that I get a bird's-eye view (US-9)
- As a learner, I want to see how concepts connect across different curricula (US-6)

## Technical Approach

- D3.js force-directed graph
- Data source: `GET /api/graph` (full graph) and `GET /api/graph/topic/:id` (per-topic subgraph) from PBI 2
- Node types:
  - Topics: large nodes with title labels
  - Concepts: small nodes (shown/hidden via filter toggle)
- Edge types (visually distinct per PRD 13.1):
  - Prerequisite: directional arrow
  - Reference: dotted line
  - Related: solid line
- Color coding (user toggle):
  - By topic membership
  - By difficulty level
  - By study progress (not started / in progress / completed)
- Interactions:
  - Click a node → navigate to that topic or concept page
  - Zoom in/out with scroll or controls
  - Zoom into a single topic → shows internal concept graph
  - Filter controls: show/hide concept nodes, show only topic-level graph
- Canvas or SVG rendering (SVG preferred for smaller graphs, Canvas for large ones)

## UX/UI Considerations

- Graph should be visually appealing — this is a motivational/exploratory feature
- Don't overwhelm with too many nodes — default to topic-level view, drill into concepts
- Legend explaining edge types and color coding
- Smooth zoom and pan transitions
- Loading state for large graphs

## Acceptance Criteria

1. Force-directed graph renders topics as large nodes and concepts as small nodes
2. Edge types visually distinct (arrows, dotted, solid per PRD)
3. Color coding toggles between topic membership, difficulty, and progress
4. Click-to-navigate works for both topic and concept nodes
5. Zoom and pan controls work (scroll + buttons)
6. Filter toggle shows/hides concept nodes
7. Per-topic subgraph view works (zoom into a single topic)
8. Legend displays for current color coding and edge types
9. Graph loads performantly with 50+ topics and 500+ concepts

## Dependencies

- **Depends on**: PBI 2 (graph API endpoints), PBI 3 (frontend scaffold)
- **External**: D3.js

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 12`._
