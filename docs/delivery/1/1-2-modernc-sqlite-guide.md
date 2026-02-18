# 1-2 modernc SQLite Guide

Date: 2026-02-18

## Package

- `modernc.org/sqlite`

## Documentation

- https://pkg.go.dev/modernc.org/sqlite
- https://gitlab.com/cznic/sqlite

## API Usage Notes

1. Register the SQL driver via blank import: `_ "modernc.org/sqlite"`.
2. Open connections with `database/sql` using driver name `sqlite`.
3. Enable pragmas per connection (for example, `PRAGMA foreign_keys = ON`).
4. SQLite URI query parameter `_pragma=...` is supported when needed.

## Apollo Usage Example

```go
package database

import (
    "context"
    "database/sql"

    _ "modernc.org/sqlite"
)

func openDatabase(ctx context.Context, path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, err
    }

    if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON;"); err != nil {
        _ = db.Close()
        return nil, err
    }

    if err := db.PingContext(ctx); err != nil {
        _ = db.Close()
        return nil, err
    }

    return db, nil
}
```

## Key Patterns

- Keep all SQL execution in `context.Context` aware paths.
- Use migration tracking (`schema_migrations`) to guarantee idempotent startup.
- Store JSON payloads as `TEXT` and validate with `json_valid(...)` constraints.
