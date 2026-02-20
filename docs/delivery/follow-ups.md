# Follow-Ups

Ideas, improvements, and deferred work captured during planning and implementation.
Review periodically — good candidates become new PBIs via `/new-pbi`.

## Open

| # | Type | Summary | Source | Date | Notes |
|---|------|---------|--------|------|-------|
| 1 | enhancement | Add progress granularity per-module during passes 2-3 | PBI-15 | 2026-02-20 | Current progress only tracks pass number. With file-per-lesson output, the orchestrator could count written lesson files to report per-module completion during passes 2 and 3. |
| 2 | perf | Parallel assembly validation per-module | PBI-15 | 2026-02-20 | The assembler could validate each module independently in parallel using goroutines, reducing assembly time for large curricula. Not needed now but useful if topics grow beyond 8 modules. |
| 3 | enhancement | Resume from partial file tree on pipeline failure | PBI-15 | 2026-02-20 | With files on disk after each pass, the orchestrator could detect which pass completed last and resume from there instead of restarting the entire pipeline. Requires tracking which files exist vs expected. |
