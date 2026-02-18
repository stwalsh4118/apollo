# PBI-9: Connection Resolver

[View in Backlog](../backlog.md)

## Overview

Implement the connection resolver that runs after each research agent returns, integrating new curricula into the existing knowledge pool. This handles concept deduplication (exact and fuzzy), cross-reference injection, prerequisite validation, and conflict detection — the system that turns isolated courses into a connected knowledge graph.

## Problem Statement

When multiple topics are researched (PBI 8), their concepts overlap. "VLAN tagging" might be defined in both a Proxmox curriculum and a Networking Fundamentals curriculum. Without the connection resolver, concepts are duplicated, cross-references are missing, and the knowledge graph is fragmented. PRD section 6.6 defines the 4-step resolution process.

## User Stories

- As a learner, I want concepts deduplicated across curricula so that each concept has one canonical definition (US-6)
- As a learner, I want cross-references between topics so that clicking a concept in one course takes me to its definition in another (US-6)

## Technical Approach

- Connection resolver runs automatically after each research job's curriculum is stored (hook into PBI 6/8 flow)
- **Step 1: Concept Deduplication**
  - Exact slug match → merge, pick more complete definition as canonical
  - Fuzzy match detection: spawn a lightweight Claude Code CLI session with the candidate pairs and existing concepts, ask "are these the same concept?" If yes → merge and alias the slug
  - No match → create new concept
- **Step 2: Cross-Reference Injection**
  - For each `concepts_referenced` in lessons, verify the concept exists
  - Exists → create bidirectional link in `concept_references` table
  - Doesn't exist → create placeholder concept with `status: 'unresolved'`
  - When an unresolved concept is later defined by another topic's research, resolve it
- **Step 3: Prerequisite Validation**
  - Lightweight check (CLI session or heuristic): are the prerequisites reasonable?
  - Flag obviously wrong prerequisites for user review
- **Step 4: Conflict Detection**
  - If two curricula define the same concept with materially different definitions → set `status: 'conflict'`
  - Surface conflicts in the API for user resolution (concept detail endpoint shows conflict status)

## UX/UI Considerations

N/A — backend PBI. Conflict surfacing in the UI is handled by existing concept detail views (PBI 3/11).

## Acceptance Criteria

1. Exact slug match deduplication works: two curricula defining the same concept slug are merged
2. Fuzzy match detection identifies near-matches (e.g., "linux-bridge" vs "linux-bridge-interface")
3. Merged concepts retain aliases in the `aliases` JSON column
4. Cross-references created bidirectionally when referenced concepts exist
5. Unresolved concepts created when referenced concepts don't exist yet
6. Unresolved concepts resolved automatically when the defining topic is later researched
7. Conflicts detected and flagged with `status: 'conflict'`
8. Prerequisite validation flags unreasonable prerequisites
9. Connection resolver runs automatically after research completion (not manual)

## Dependencies

- **Depends on**: PBI 8 (orchestrator produces multiple curricula that need resolution), PBI 6 (single topic pipeline)
- **External**: Claude Code CLI (for fuzzy matching sessions)

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 9`._
