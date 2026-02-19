# 3-2 TanStack Query (React Query v5) Guide

Date: 2026-02-18

## Package

- `@tanstack/react-query` (v5)

## Documentation

- https://tanstack.com/query/v5/docs/framework/react/overview

---

## 1. Installation

```bash
pnpm add @tanstack/react-query
```

Optional dev tools (useful during development, tree-shaken in production):

```bash
pnpm add -D @tanstack/react-query-devtools
```

Requires React 18+.

---

## 2. QueryClientProvider Setup

Every application using TanStack Query must wrap its component tree with `QueryClientProvider`. Create the `QueryClient` instance **outside** the component to avoid recreating it on every render.

```tsx
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,   // 5 minutes — data considered fresh
      gcTime: 10 * 60 * 1000,     // 10 minutes — unused cache kept in memory
      retry: 1,                    // retry failed requests once
      refetchOnWindowFocus: false, // disable automatic refetch on tab focus
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      {/* application routes / components */}
    </QueryClientProvider>
  );
}
```

---

## 3. Typed Query Hooks with useQuery

### Basic Typed Hook

TanStack Query v5 infers types from the `queryFn` return type. You do **not** need to pass generic type parameters manually if your fetch function has a typed return.

```tsx
import { useQuery } from "@tanstack/react-query";

interface Topic {
  id: number;
  title: string;
  difficulty: string;
}

async function fetchTopics(): Promise<Topic[]> {
  const res = await fetch("/api/v1/topics");
  if (!res.ok) {
    throw new Error(`Failed to fetch topics: ${res.status}`);
  }
  return res.json();
}

// `data` is inferred as `Topic[] | undefined`
export function useTopics() {
  return useQuery({
    queryKey: ["topics"],
    queryFn: fetchTopics,
  });
}
```

### Hook with Parameters

```tsx
interface TopicDetail {
  id: number;
  title: string;
  modules: ModuleSummary[];
}

async function fetchTopicDetail(id: number): Promise<TopicDetail> {
  const res = await fetch(`/api/v1/topics/${id}`);
  if (!res.ok) {
    throw new Error(`Failed to fetch topic ${id}: ${res.status}`);
  }
  return res.json();
}

export function useTopicDetail(id: number) {
  return useQuery({
    queryKey: ["topics", id, "detail"],
    queryFn: () => fetchTopicDetail(id),
    enabled: id > 0, // only run when id is valid
  });
}
```

### Using the queryOptions Helper

The `queryOptions` helper co-locates `queryKey` and `queryFn` in a single object with full type inference. This is the recommended pattern when the same query config is used in multiple places (hooks, prefetching, cache invalidation).

```tsx
import { queryOptions, useQuery } from "@tanstack/react-query";

function topicDetailOptions(id: number) {
  return queryOptions({
    queryKey: ["topics", id, "detail"] as const,
    queryFn: () => fetchTopicDetail(id),
    enabled: id > 0,
  });
}

// In a component:
export function useTopicDetail(id: number) {
  return useQuery(topicDetailOptions(id));
}

// For prefetching:
queryClient.prefetchQuery(topicDetailOptions(id));

// For cache invalidation:
queryClient.invalidateQueries({ queryKey: ["topics", id, "detail"] });
```

---

## 4. queryKey Conventions

Keys must be arrays. Structure them **from most generic to most specific**.

### Query Key Factory Pattern

Define one factory object per feature/entity:

```tsx
export const topicKeys = {
  all:     ["topics"] as const,
  lists:   () => [...topicKeys.all, "list"] as const,
  list:    (filters: string) => [...topicKeys.lists(), { filters }] as const,
  details: () => [...topicKeys.all, "detail"] as const,
  detail:  (id: number) => [...topicKeys.details(), id] as const,
  full:    (id: number) => [...topicKeys.all, "full", id] as const,
};

export const moduleKeys = {
  all:     ["modules"] as const,
  details: () => [...moduleKeys.all, "detail"] as const,
  detail:  (id: number) => [...moduleKeys.details(), id] as const,
};

export const lessonKeys = {
  all:     ["lessons"] as const,
  details: () => [...lessonKeys.all, "detail"] as const,
  detail:  (id: number) => [...lessonKeys.details(), id] as const,
};
```

### Invalidation Granularity

The hierarchical structure enables targeted cache invalidation:

```tsx
// Invalidate everything related to topics:
queryClient.invalidateQueries({ queryKey: topicKeys.all });

// Invalidate only topic list queries:
queryClient.invalidateQueries({ queryKey: topicKeys.lists() });

// Invalidate a single topic detail:
queryClient.invalidateQueries({ queryKey: topicKeys.detail(42) });
```

Keys are matched using a **prefix** strategy: `["topics"]` matches `["topics", "list"]`, `["topics", "detail", 42]`, etc.

---

## 5. Configuration Options

### Important Defaults (v5)

| Option | Default | Description |
|--------|---------|-------------|
| `staleTime` | `0` | Data is stale immediately after fetch |
| `gcTime` | `300000` (5 min) | Unused cache entries garbage-collected after this |
| `retry` | `3` | Failed queries retry 3 times with exponential backoff |
| `refetchOnWindowFocus` | `true` | Refetch stale queries when window regains focus |
| `refetchOnMount` | `true` | Refetch stale queries when component mounts |
| `refetchOnReconnect` | `true` | Refetch stale queries when network reconnects |

### Recommended Overrides for Apollo

```tsx
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,   // curriculum data changes infrequently
      gcTime: 10 * 60 * 1000,     // keep cache for 10 minutes
      retry: 1,                    // single retry on failure
      refetchOnWindowFocus: false, // avoid unnecessary refetches for static content
    },
  },
});
```

**Rule of thumb**: `gcTime` should always be >= `staleTime`. If cache is garbage-collected before data goes stale, a new fetch is forced anyway.

### Per-Query Overrides

Any default can be overridden at the individual query level:

```tsx
useQuery({
  queryKey: topicKeys.detail(id),
  queryFn: () => fetchTopicDetail(id),
  staleTime: Infinity, // never mark this data as stale
  gcTime: 30 * 60 * 1000, // keep for 30 minutes
});
```

---

## 6. Error Handling Patterns

### Pattern A: Component-Level Error State

The `useQuery` return object provides `error` and `isError` for local handling:

```tsx
function TopicList() {
  const { data, isLoading, isError, error } = useTopics();

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error: {error.message}</div>;

  return (
    <ul>
      {data.map((topic) => (
        <li key={topic.id}>{topic.title}</li>
      ))}
    </ul>
  );
}
```

### Pattern B: Error Boundaries with throwOnError

Propagate errors to the nearest React error boundary. Accepts a boolean or a function for selective propagation:

```tsx
// Throw all errors to boundary
useQuery({
  queryKey: topicKeys.detail(id),
  queryFn: () => fetchTopicDetail(id),
  throwOnError: true,
});

// Only throw server errors (5xx), handle 4xx locally
useQuery({
  queryKey: topicKeys.detail(id),
  queryFn: () => fetchTopicDetail(id),
  throwOnError: (error) => {
    return error instanceof ApiError && error.status >= 500;
  },
});
```

Use `QueryErrorResetBoundary` to reset and retry from error boundaries:

```tsx
import { QueryErrorResetBoundary } from "@tanstack/react-query";
import { ErrorBoundary } from "react-error-boundary";

function TopicPage() {
  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ErrorBoundary onReset={reset} fallbackRender={({ resetErrorBoundary }) => (
          <div>
            <p>Something went wrong.</p>
            <button onClick={resetErrorBoundary}>Retry</button>
          </div>
        )}>
          <TopicDetail />
        </ErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  );
}
```

### Pattern C: Global Error Callbacks

In v5, `onError`/`onSuccess` callbacks are removed from `useQuery`. Use `QueryCache`-level callbacks for global handling (toasts, logging):

```tsx
import { QueryClient, QueryCache } from "@tanstack/react-query";

const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, query) => {
      // Only toast for background refetch failures (stale data exists)
      if (query.state.data !== undefined) {
        toast.error(`Background update failed: ${error.message}`);
      }
    },
  }),
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
    },
  },
});
```

### Custom Error Type

Define a typed error class for API errors to enable conditional handling:

```tsx
export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function fetchJson<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) {
    throw new ApiError(res.status, `${res.statusText}: ${url}`);
  }
  return res.json();
}
```

---

## Apollo Usage Summary

For task 3-2, the key integration points are:

1. **Install**: `pnpm add @tanstack/react-query`
2. **Types**: Define in `web/src/api/types.ts` (no TanStack Query dependency)
3. **Client**: Fetch functions in `web/src/api/client.ts` returning typed promises
4. **Hooks**: TanStack Query hooks in `web/src/api/hooks.ts` using `queryKey` factories
5. **Provider**: Wrap app in `QueryClientProvider` in `web/src/main.tsx`
6. **Errors**: Use component-level `isError`/`error` pattern as default; add global `QueryCache.onError` for toast notifications if needed later
