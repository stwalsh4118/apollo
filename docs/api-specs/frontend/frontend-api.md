# Frontend API Specification

## Overview

React + TypeScript SPA served by Vite, consuming the Go curriculum API via `/api` proxy.

## Routes

| Path | Page | Description |
|------|------|-------------|
| `/` | TopicListPage | Topic cards grid with title, difficulty, description, module count |
| `/topics/:id` | CourseViewPage | Module sidebar + lesson content with all 6 section renderers |
| `*` | NotFoundPage | 404 catch-all |

## Component Hierarchy

```
App (RouterProvider)
└── AppLayout (header + Outlet)
    ├── TopicListPage
    │   └── TopicCard + DifficultyBadge
    ├── CourseViewPage
    │   ├── ModuleSidebar
    │   │   └── ModuleItem → LessonItem
    │   ├── LessonContent
    │   │   ├── ContentRenderer
    │   │   │   ├── TextSection (react-markdown + rehype-raw)
    │   │   │   ├── CodeSection (Shiki highlighting)
    │   │   │   ├── CalloutSection (info/tip/warning/prerequisite)
    │   │   │   ├── DiagramSection (Mermaid.js)
    │   │   │   ├── TableSection
    │   │   │   ├── ImageSection
    │   │   │   └── UnknownSection (fallback)
    │   │   ├── ExerciseList → ExerciseBlock
    │   │   └── ReviewQuestions
    │   └── LessonNavigation (prev/next)
    └── NotFoundPage
```

## API Client

### Base Configuration

- Proxy: Vite dev server forwards `/api/*` → `http://localhost:8080`
- Error class: `ApiError { status: number; message: string }`

### Client Functions

```typescript
function fetchJson<T>(url: string): Promise<T>
function fetchTopics(): Promise<PaginatedResponse<TopicSummary>>
function fetchTopicDetail(id: string): Promise<TopicDetail>
function fetchTopicFull(id: string): Promise<TopicFull>
function fetchModuleDetail(id: string): Promise<ModuleDetail>
function fetchLessonDetail(id: string): Promise<LessonDetail>
```

## TanStack Query Hooks

### Query Key Factories

```typescript
const topicKeys = { all: ["topics"], detail: (id) => ["topics", id], full: (id) => ["topics", id, "full"] }
const moduleKeys = { all: ["modules"], detail: (id) => ["modules", id] }
const lessonKeys = { all: ["lessons"], detail: (id) => ["lessons", id] }
```

### Hooks

```typescript
function useTopics(): UseQueryResult<PaginatedResponse<TopicSummary>>
function useTopicDetail(id: string): UseQueryResult<TopicDetail>
function useTopicFull(id: string): UseQueryResult<TopicFull>
function useModuleDetail(id: string): UseQueryResult<ModuleDetail>
function useLessonDetail(id: string): UseQueryResult<LessonDetail>
```

### Query Client Config

| Setting | Value |
|---------|-------|
| staleTime | 5 min |
| gcTime | 10 min |
| retry | 1 |
| refetchOnWindowFocus | false |

## TypeScript Types

### Content Sections (discriminated union on `type`)

```typescript
type ContentSection = TextSection | CodeSection | CalloutSection | DiagramSection | TableSection | ImageSection

interface TextSection { type: "text"; body: string }
interface CodeSection { type: "code"; language: string; code: string; title?: string; explanation?: string }
interface CalloutSection { type: "callout"; variant: "prerequisite" | "warning" | "tip" | "info"; body: string; concept_ref?: string }
interface DiagramSection { type: "diagram"; format: "mermaid" | "image"; source: string; title?: string; url?: string; alt?: string; caption?: string }
interface TableSection { type: "table"; headers: string[]; rows: string[][] }
interface ImageSection { type: "image"; url: string; alt: string; caption?: string }
```

### Entity Types

```typescript
// Topics
interface TopicSummary { id, title, description, difficulty, status, estimated_hours, tags, module_count, created_at, updated_at }
interface TopicDetail extends TopicSummary { modules: ModuleSummary[] }
interface TopicFull extends TopicSummary { modules: ModuleFull[] }

// Modules
interface ModuleSummary { id, topic_id, title, description, sort_order, estimated_minutes, lesson_count, created_at, updated_at }
interface ModuleDetail extends ModuleSummary { lessons: LessonSummary[] }
interface ModuleFull extends ModuleSummary { lessons: LessonFull[] }

// Lessons
interface LessonSummary { id, module_id, title, sort_order, estimated_minutes, created_at, updated_at }
interface LessonDetail extends LessonSummary { content: ContentSection[]; examples: Example[]; exercises: Exercise[]; review_questions: ReviewQuestion[] }

// Supporting
interface Exercise { type: string; title: string; instructions: string; hints: string[]; success_criteria: string[]; environment?: string }
interface ReviewQuestion { question: string; answer: string; concepts_tested: string[] }
```

## Custom Hooks

### useShikiHighlighter

```typescript
function useShikiHighlighter(code: string, language: string): { html: string | null; isLoading: boolean }
```

- Lazy singleton highlighter (github-dark theme)
- Languages: bash, json, yaml, go, javascript, typescript
- Derives loading state from result staleness

### useMermaid

```typescript
function useMermaid(source: string): { svg: string | null; error: string | null; isLoading: boolean }
```

- Lazy dynamic import, strict security, dark theme
- Unique render IDs per invocation

## Utilities

### lessonNavigation

```typescript
function getLessonNavigation(modules: ModuleFull[], currentLessonId: string): { prev: LessonNavTarget | null; next: LessonNavTarget | null }
```

Flattens modules into linear lesson list and returns adjacent entries.
