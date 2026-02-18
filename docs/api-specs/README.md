# API Specifications Index

This directory contains concise API specifications to prevent code duplication
and ensure consistency across the project.

## Available Specs

| System | Spec File | Description |
|--------|-----------|-------------|
| Database | [database-api.md](./database/database-api.md) | SQLite connection, migration system, full schema reference |
| Curriculum | [curriculum-api.md](./curriculum/curriculum-api.md) | REST endpoints, repository interfaces, models, error handling |

## Organization

- One directory per system/domain: `docs/api-specs/<system>/`
- One markdown file per logical grouping: `<system>-api.md`
- See `~/.claude/_references/api-specs.md` for format guidelines.
