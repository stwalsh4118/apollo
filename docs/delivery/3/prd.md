# PBI-3: Frontend Foundation & Course View

[View in Backlog](../backlog.md)

## Overview

Scaffold the React + TypeScript frontend and build the primary learning interface: the course view. This is the first visual surface of Apollo — a user can browse topics, navigate modules via a sidebar, and read fully rendered lessons with all 6 content section types (text, code, callout, diagram, table, image).

## Problem Statement

The API serves curriculum data (PBI 2) but there's no way to see or interact with it. The course view is the primary learning experience — it needs to render structured lesson content with syntax highlighting, mermaid diagrams, callout boxes, and inline examples. This is how the user will spend most of their time in Apollo.

## User Stories

- As a learner, I want to see all my topics as browsable cards so that I can pick what to study
- As a learner, I want a module sidebar showing my place in a curriculum so that I can navigate between lessons
- As a learner, I want lesson content rendered with syntax highlighting, diagrams, and styled callouts so that the material is readable and engaging
- As a learner, I want to see exercises and review questions within lessons so that I can practice

## Technical Approach

- React + TypeScript + Vite
- Tailwind CSS for styling
- React Router v7 for client-side routing
- TanStack Query for data fetching and caching
- Routes: `/` (topic list), `/topics/:id` (course view with module sidebar + lesson)
- Content section renderers:
  - `text` → react-markdown with rehype plugins
  - `code` → Shiki for syntax highlighting
  - `callout` → styled callout box (prerequisite/warning/tip/info variants)
  - `diagram` → Mermaid.js renderer
  - `table` → HTML table from headers + rows
  - `image` → `<img>` with caption
- Module sidebar: collapsible module list, lesson links, completion indicators (wired in PBI 4)
- Exercise rendering: instructions, hints (progressive reveal), success criteria
- Review questions: collapsible section at lesson end

## UX/UI Considerations

- Course view is the primary screen — clean, readable, distraction-free
- Module sidebar on the left, lesson content in the main area (responsive for smaller screens)
- Code blocks need copy button and language label
- Mermaid diagrams render inline with fallback for parse errors
- Callout boxes visually distinct by variant (color/icon)
- Exercises and review questions in collapsible sections so they don't overwhelm the lesson flow

## Acceptance Criteria

1. React app builds and serves via Vite dev server
2. Topic list page shows topic cards with title, difficulty, description, and module count
3. Course view renders module sidebar with all modules and their lessons
4. All 6 content section types render correctly (text, code, callout, diagram, table, image)
5. Shiki syntax highlighting works for at least bash, json, yaml, go, javascript, typescript
6. Mermaid diagrams render from `source` field with error fallback
7. Exercises render with progressive hint reveal
8. Review questions render in a collapsible section
9. Navigation between lessons works (prev/next within module, cross-module)
10. API proxy configured so frontend dev server can reach the Go API

## Dependencies

- **Depends on**: PBI 2 (API must serve curriculum data)
- **External**: React, Vite, Tailwind, React Router v7, TanStack Query, Shiki, Mermaid.js, react-markdown

## Open Questions

- None

## Related Tasks

[View Tasks](./tasks.md)
