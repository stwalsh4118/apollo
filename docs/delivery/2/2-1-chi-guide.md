# Chi Router v5 — API Guide

**Date**: 2026-02-18
**Docs**: https://pkg.go.dev/github.com/go-chi/chi/v5
**Version**: v5.2.5

## Router Creation

```go
r := chi.NewRouter()
```

## Middleware

```go
r.Use(middleware.Logger)      // Request logging (place before Recoverer)
r.Use(middleware.Recoverer)   // Panic recovery → 500
r.Use(middleware.SetHeader("Content-Type", "application/json"))
```

## Routes

```go
r.Get("/path", handler)
r.Post("/path", handler)
r.Put("/path", handler)
r.Delete("/path", handler)
```

## URL Parameters

Pattern: `{name}` in route path.

```go
r.Get("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
})
```

## Route Grouping

```go
r.Route("/api/topics", func(r chi.Router) {
    r.Get("/", listTopics)
    r.Route("/{id}", func(r chi.Router) {
        r.Get("/", getTopicByID)
        r.Get("/full", getTopicFull)
    })
})
```

## Custom Middleware Pattern

```go
func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // before
        next.ServeHTTP(w, r)
        // after
    })
}
```

## Key Middleware

| Middleware | Function |
|-----------|----------|
| `middleware.Logger` | Log requests with duration |
| `middleware.Recoverer` | Recover panics → 500 |
| `middleware.SetHeader(k,v)` | Set response header |
| `middleware.RequestID` | Inject request ID |
| `middleware.RealIP` | Extract real client IP |
| `middleware.Timeout(d)` | Context timeout |
