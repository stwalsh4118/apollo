# PBI-5: Research Skill Prompt & Schema

[View in Backlog](../backlog.md)

## Overview

Create the research skill prompt file that drives Claude Code research sessions and the curriculum JSON schema used for structured output validation. These are the "instructions" and "contract" that make the research pipeline produce consistent, high-quality output.

## Problem Statement

The research agent (PBI 6) needs two things before it can run: (1) a detailed skill prompt that implements the 4-pass research pipeline (survey → deep dive → exercises → self-review) with the self-review checklist, and (2) a JSON schema that enforces the curriculum structure from PRD section 7. Without these, Claude Code sessions have no instructions and no output contract.

## User Stories

- As a developer, I want a research skill prompt that guides Claude Code through a structured 4-pass research pipeline so that output quality is consistent
- As a developer, I want a JSON schema matching the curriculum spec so that structured output is validated automatically via `--json-schema`

## Technical Approach

- Research skill prompt file (`skills/research.md` or similar):
  - Pass 1 (Survey): broad web search, module structure, prerequisite identification, split detection
  - Pass 2 (Deep Dive): per-module focused research, concept identification, content sections
  - Pass 3 (Exercises): exercise generation (all 7 types from PRD 6.2), review questions, assessments
  - Pass 4 (Self-Review): checklist validation per PRD 6.1
  - Knowledge pool context integration instructions (read `knowledge_pool_summary.json`)
  - Prerequisite classification rules (essential/helpful/deep_background per PRD 5.3)
  - Topic splitting rules (>8 modules per PRD 5.4)
- Curriculum JSON schema (`schemas/curriculum.json`):
  - Matches PRD section 7 exactly: Topic → Modules → Lessons → Concepts
  - All content section types (text, code, callout, diagram, table, image)
  - Exercise types with required fields
  - Prerequisite classification structure
  - Used with `--json-schema` flag on the final research pass
- Go-side schema validation:
  - Load and validate returned JSON against the schema before ingesting
  - Clear error messages for validation failures

## UX/UI Considerations

N/A — developer tooling / infrastructure PBI.

## Acceptance Criteria

1. Research skill prompt file covers all 4 passes with clear instructions per pass
2. Self-review checklist from PRD 6.1 is included in the prompt
3. Exercise type spectrum from PRD 6.2 is documented in the prompt with usage guidance
4. Prerequisite classification rules from PRD 5.3 are included
5. Topic splitting rules from PRD 5.4 are included
6. JSON schema validates a correct curriculum object (positive test)
7. JSON schema rejects malformed objects — missing required fields, wrong types (negative tests)
8. Go validation function loads the schema and validates JSON input
9. Knowledge pool summary format documented and schema created

## Dependencies

- **Depends on**: PBI 1 (Go project must exist for validation code)
- **External**: JSON Schema spec, Claude Code `--json-schema` flag behaviour

## Open Questions

- None

## Related Tasks

[View Tasks](./tasks.md)
