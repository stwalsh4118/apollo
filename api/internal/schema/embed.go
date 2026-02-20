package schema

import "embed"

// schemaFS embeds JSON schema files for validation at runtime.
// Canonical copies live at schemas/ (project root) for use with
// Claude Code's --json-schema flag. These embedded copies are kept
// in sync and used by the Go validation logic.
//
//go:embed curriculum.json knowledge_pool_summary.json
var schemaFS embed.FS

// CurriculumSchemaJSON returns the raw bytes of the embedded curriculum JSON schema.
// Used by the research orchestrator to write the schema to the job work directory.
func CurriculumSchemaJSON() ([]byte, error) {
	return schemaFS.ReadFile(curriculumSchemaFile)
}
