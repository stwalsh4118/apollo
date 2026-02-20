# PBI-15: Research Pipeline — File-Per-Lesson Output Structure

[View in Backlog](../backlog.md)

## Overview

Restructure the research pipeline's intermediate output from a single monolithic `curriculum.json` (3K+ lines) to a hierarchical file-per-lesson directory tree. Sub-agents write individual files, each pass operates on small per-file units, and Go code handles final assembly and ingestion. This keeps every file under ~300 lines — comfortably within any agent's context window.

## Problem Statement

The current pipeline produces a single `curriculum.json` that exceeds 3,000 lines for a typical topic. This monolithic blob:

- **Cannot fit in a single agent context window**, making assembly and validation passes unreliable
- **Requires an assembly bottleneck** where one agent must hold all sub-agent outputs in memory
- **Makes partial progress invisible** — if the pipeline fails mid-way, nothing is recoverable
- **Prevents per-lesson parallelization** in passes 3 and 4

A file-per-lesson directory structure solves all of these: each agent reads/writes only the files it needs, Go code merges the tree, and partial results survive failures.

## User Stories

- As a developer, I want the research pipeline to produce small, composable files so that each agent works within its context window and the pipeline is reliable

## Technical Approach

### Directory Structure

Each research job produces a file tree in its work directory:

```
data/research/<job-id>/
  knowledge_pool_summary.json   # (existing) context for the agent
  research.md                   # (existing) embedded system prompt
  curriculum.json               # (existing) embedded schema
  topic.json                    # Pass 1 output: metadata, prerequisites, module plan
  modules/
    01-<module-slug>/
      module.json               # Module metadata: title, description, objectives, assessment
      01-<lesson-slug>.json     # Lesson content, concepts_taught, concepts_referenced
      02-<lesson-slug>.json
    02-<module-slug>/
      module.json
      01-<lesson-slug>.json
      ...
```

### Pass Changes

**Pass 1 (Survey):** Agent writes `topic.json` containing:
- Topic-level fields (id, title, description, difficulty, tags, prerequisites, related_topics, source_urls)
- Module plan: array of `{id, title, description, order}` stubs (no lessons yet)
- Creates `modules/<NN>-<slug>/` directories

**Pass 2 (Deep Dive):** Main agent does research, then sub-agents write directly:
- Each sub-agent writes 1-2 lesson JSON files to `modules/<slug>/<NN>-<lesson>.json`
- Each sub-agent writes/updates `modules/<slug>/module.json` with learning objectives
- No assembly step — files land directly in the tree

**Pass 3 (Exercises):** Sub-agents read existing lesson files, add exercises and review questions, write back:
- Read `modules/<slug>/<lesson>.json`, add `exercises` and `review_questions` fields, write back
- Each sub-agent also writes the module `assessment` into `module.json`

**Pass 4 (Validation):** Agent reads the file tree, validates per-file quality, fixes issues:
- No `--json-schema` needed — validation happens in Go after assembly
- Agent focuses on content quality checks (the self-review checklist)

### Go Assembler

New Go code that:
1. Walks the `modules/` directory tree
2. Reads `topic.json` + all `module.json` + all lesson files
3. Assembles into `CurriculumOutput` struct (existing type, unchanged)
4. Validates against curriculum schema
5. Passes to existing `CurriculumIngester` for SQLite storage

### Orchestrator Changes

- Pass 4 no longer uses `--json-schema` flag
- After pass 4 completes, orchestrator calls the Go assembler instead of extracting `structured_output`
- `runFinalPass` replaced by `runPass` + `assembleFromDir`

### System Prompt Changes

The research prompt is updated to instruct agents to:
- Write files using the Write tool instead of accumulating content in memory
- Follow the directory naming convention (`NN-slug` for ordering)
- Read/modify existing files when adding exercises (pass 3)

## UX/UI Considerations

N/A — internal pipeline change. API responses unchanged.

## Acceptance Criteria

1. Pass 1 writes `topic.json` and creates module directories
2. Pass 2 sub-agents write individual lesson files (~100-200 lines each) to `modules/<slug>/`
3. Pass 3 sub-agents read lesson files, add exercises/review questions, write back
4. No individual file in the work directory exceeds ~300 lines
5. Go assembler reads the file tree and produces a valid `CurriculumOutput`
6. Assembled output passes curriculum schema validation
7. Existing ingestion stores the same data in SQLite (API responses unchanged)
8. Partial progress survives — if pass 3 fails, pass 2 lesson files are still on disk
9. System prompt updated with file-writing instructions and directory conventions

## Dependencies

- **Depends on**: PBI 6 (research pipeline foundation)
- **Blocks**: None (PBI 7/8 can proceed in parallel)
- **External**: None

## Open Questions

- None

## Related Tasks

[View Tasks](./tasks.md)
