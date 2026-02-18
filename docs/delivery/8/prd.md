# PBI-8: Research Orchestrator & Queue

[View in Backlog](../backlog.md)

## Overview

Extend the single-topic research pipeline (PBI 6) into the full recursive orchestrator described in PRD section 5.2. The orchestrator manages a research queue, auto-expands essential prerequisites up to depth 3, runs parallel research agents (configurable max), and tracks the expansion queue for helpful/deep topics.

## Problem Statement

PBI 6 handles one topic at a time. The PRD's core value proposition is recursive research — submit "Proxmox" and automatically get Linux Administration and Networking Fundamentals researched too. The orchestrator needs to extract prerequisites from completed research, check the knowledge pool for what already exists, queue missing essentials, respect depth limits, and run multiple agents in parallel.

## User Stories

- As a learner, I want prerequisite topics auto-generated alongside my main topic so that I fill knowledge gaps (US-2)
- As a learner, I want to expand "helpful" or "deep background" topics on demand (US-8)

## Technical Approach

- After a topic's research completes and is stored (PBI 6 flow), the orchestrator:
  1. Reads `prerequisites.essential` from the curriculum
  2. For each essential: checks if it exists in the knowledge pool
  3. Missing essentials → add to research queue at `depth + 1` (if depth < MAX_RESEARCH_DEPTH)
  4. Helpful/deep → store in `expansion_queue` as `available`
- Research queue processing:
  - FIFO with priority weighting (lower depth = higher priority)
  - Configurable `MAX_PARALLEL_AGENTS` (default 3) — run up to N research sessions concurrently
  - Each queued item becomes a research job (PBI 6 pipeline)
- Expansion API:
  - `POST /api/research/expand/:topicId` — move a helpful/deep topic from `available` to `queued`
  - `POST /api/research/refresh/:topicId` — re-research an existing topic (PRD section 12)
- Knowledge pool checks prevent duplicate research (topic already exists → skip)
- Root job tracks all spawned sub-jobs for progress reporting

## UX/UI Considerations

N/A — backend PBI. Frontend expansion triggers come in PBI 13 (Dashboard).

## Acceptance Criteria

1. After a topic is published, its essential prerequisites are checked against the knowledge pool
2. Missing essentials are automatically queued for research at depth + 1
3. Depth limit (default 3) prevents infinite recursion
4. Parallel execution: up to MAX_PARALLEL_AGENTS research sessions run concurrently
5. Helpful and deep_background prerequisites stored in expansion_queue as `available`
6. `POST /api/research/expand/:topicId` queues an available expansion topic for research
7. `POST /api/research/refresh/:topicId` starts a refresh job with existing curriculum as context
8. Topics already in the knowledge pool are not re-researched (skipped in queue)
9. Research job list shows both root and prerequisite jobs with their relationships

## Dependencies

- **Depends on**: PBI 6 (single-topic pipeline must work first)
- **External**: None

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 8`._
