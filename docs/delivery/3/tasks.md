# Tasks for PBI 3: Frontend Foundation & Course View

This document lists all tasks associated with PBI 3.

**Parent PBI**: [PBI 3: Frontend Foundation & Course View](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 3-1 | [React + Vite + Tailwind scaffold with API proxy](./3-1.md) | Proposed | Scaffold React/TypeScript/Vite project, configure Tailwind CSS, set up dev server API proxy to Go backend |
| 3-2 | [TypeScript types & TanStack Query API client](./3-2.md) | Proposed | Define TypeScript interfaces matching Go API models; create TanStack Query hooks for all curriculum endpoints |
| 3-3 | [App routing & layout shell](./3-3.md) | Proposed | Configure React Router v7 routes and build responsive app layout with sidebar/main content areas |
| 3-4 | [Topic list page](./3-4.md) | Proposed | Build the `/` route with topic cards showing title, difficulty, description, and module count |
| 3-5 | [Course view layout & module sidebar](./3-5.md) | Proposed | Build the `/topics/:id` course view with collapsible module sidebar and lesson content area |
| 3-6 | [Content section renderers — text, table, callout, image](./3-6.md) | Proposed | Implement four content section renderers: text (react-markdown), table, callout (4 variants), and image |
| 3-7 | [Code section renderer with Shiki](./3-7.md) | Proposed | Implement code section renderer with Shiki syntax highlighting, copy button, and language label |
| 3-8 | [Diagram renderer with Mermaid.js](./3-8.md) | Proposed | Implement diagram section renderer with Mermaid.js rendering and error fallback |
| 3-9 | [Exercise & review question rendering](./3-9.md) | Proposed | Render exercises with progressive hint reveal and review questions in collapsible sections |
| 3-10 | [Lesson navigation](./3-10.md) | Proposed | Implement prev/next lesson navigation within and across modules |
| 3-11 | [E2E CoS test](./3-11.md) | Proposed | End-to-end verification of all PBI 3 acceptance criteria |

## Dependency Graph

```
3-1  (Scaffold)
├── 3-2  (Types & API Client)
│   ├── 3-4  (Topic List Page)
│   └── 3-5  (Course View & Sidebar)
│       ├── 3-6  (Text, Table, Callout, Image renderers)
│       │   ├── 3-7  (Code renderer + Shiki)
│       │   └── 3-8  (Diagram renderer + Mermaid)
│       ├── 3-9  (Exercises & Review Questions)
│       └── 3-10 (Lesson Navigation)
├── 3-3  (Routing & Layout)
│   ├── 3-4  (Topic List Page)
│   └── 3-5  (Course View & Sidebar)
└───────────────────────────────────────
    3-11 (E2E CoS Test) ← depends on ALL above
```

## Implementation Order

1. **3-1** — React + Vite + Tailwind scaffold. Foundation for everything.
2. **3-2** — TypeScript types & API client. Required by all data-fetching components.
3. **3-3** — Routing & layout shell. Required by all pages.
4. **3-4** — Topic list page. First visible feature; depends on 3-2 and 3-3.
5. **3-5** — Course view & module sidebar. Core learning UI; depends on 3-2 and 3-3.
6. **3-6** — Text, table, callout, image renderers. First content rendering; depends on 3-5.
7. **3-7** — Code renderer with Shiki. Extends content rendering; depends on 3-6.
8. **3-8** — Diagram renderer with Mermaid. Extends content rendering; depends on 3-6.
9. **3-9** — Exercise & review question rendering. Depends on 3-5 (lesson content area).
10. **3-10** — Lesson navigation. Depends on 3-5 (needs topic tree and sidebar).
11. **3-11** — E2E CoS test. Must be last; depends on all tasks above.

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 3-1 | Simple | Vite, Tailwind CSS v4 |
| 3-2 | Medium | @tanstack/react-query |
| 3-3 | Simple | react-router (v7) |
| 3-4 | Simple | — |
| 3-5 | Medium | — |
| 3-6 | Medium | react-markdown, rehype-raw |
| 3-7 | Complex | shiki |
| 3-8 | Medium | mermaid |
| 3-9 | Medium | — |
| 3-10 | Medium | — |
| 3-11 | Complex | Playwright (MCP) |

## External Package Research Required

| Task ID | Package | Guide Document |
|---------|---------|----------------|
| 3-2 | @tanstack/react-query | `3-2-tanstack-query-guide.md` |
| 3-3 | react-router v7 | `3-3-react-router-guide.md` |
| 3-6 | react-markdown | `3-6-react-markdown-guide.md` |
| 3-7 | shiki | `3-7-shiki-guide.md` |
| 3-8 | mermaid | `3-8-mermaid-guide.md` |
