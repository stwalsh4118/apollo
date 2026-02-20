# Follow-Ups

Ideas, improvements, and deferred work captured during planning and implementation.
Review periodically — good candidates become new PBIs via `/new-pbi`.

## Open

| # | Type | Summary | Source | Date | Notes |
|---|------|---------|--------|------|-------|
| 1 | enhancement | Add progress granularity per-module during passes 2-3 | PBI-15 | 2026-02-20 | Current progress only tracks pass number. With file-per-lesson output, the orchestrator could count written lesson files to report per-module completion during passes 2 and 3. |
| 2 | perf | Parallel assembly validation per-module | PBI-15 | 2026-02-20 | The assembler could validate each module independently in parallel using goroutines, reducing assembly time for large curricula. Not needed now but useful if topics grow beyond 8 modules. |
| 3 | enhancement | Resume from partial file tree on pipeline failure | PBI-15 | 2026-02-20 | With files on disk after each pass, the orchestrator could detect which pass completed last and resume from there instead of restarting the entire pipeline. Requires tracking which files exist vs expected. |
| 4 | tech-debt | Remove dead `JSONSchemaFile` field from `ResumePassOpts` | PBI-15 | 2026-02-20 | Field is no longer set by any caller after PBI-15 removed `runFinalPass`. `BuildResumeArgs` still conditionally appends `--json-schema` and `TestBuildResumeArgsFinalPass` tests that path, but nothing exercises it in production. |
| 5 | perf | Eliminate triple serialization in assembly-to-ingestion path | PBI-15 | 2026-02-20 | `AssembleFromDir` marshals for schema validation, orchestrator marshals again for `json.RawMessage`, and `Ingest` validates schema + unmarshals. Consider `IngestFromStruct` or returning validated JSON from assembler. Low priority given infrequent execution. |
| 6 | fix | Frontend LessonContent type mismatch with schema | PBI-15 | 2026-02-20 | Frontend `types.ts` had `content: ContentSection[]` (flat array) but schema defines `LessonContent` as `{sections: ContentSection[]}`. Fixed during PBI-15 completion but indicates need for type generation or shared contract. |
| 7 | enhancement | Render markdown in exercise success_criteria and environment fields | PBI-15 | 2026-02-20 | Only `instructions` and `hints` fields currently use react-markdown. Other exercise text fields could also contain code fences or formatting. |
