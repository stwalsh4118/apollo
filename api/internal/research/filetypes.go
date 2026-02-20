package research

import "encoding/json"

// Directory and file name constants for the file-per-lesson output structure.
// Agents write individual files during passes 1-3; the assembler (AssembleFromDir)
// reads them back to produce a CurriculumOutput.
const (
	// ModulesDirName is the subdirectory under the work directory containing module directories.
	ModulesDirName = "modules"

	// ModuleFileBaseName is the metadata file inside each module directory.
	ModuleFileBaseName = "module.json"

	// TopicFileName is the topic metadata file in the work directory root.
	TopicFileName = "topic.json"
)

// TopicFile represents the topic.json written by Pass 1 (Survey).
// It contains topic-level metadata, prerequisites, and a module plan
// with lightweight stubs (no lessons yet).
type TopicFile struct {
	ID             string              `json:"id"`
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	Difficulty     string              `json:"difficulty"`
	EstimatedHours float64             `json:"estimated_hours"`
	Tags           []string            `json:"tags"`
	Prerequisites  PrerequisitesOutput `json:"prerequisites"`
	RelatedTopics  []string            `json:"related_topics"`
	SourceURLs     []string            `json:"source_urls"`
	GeneratedAt    string              `json:"generated_at"`
	Version        int                 `json:"version"`
	ModulePlan     []ModulePlanEntry   `json:"module_plan"`
}

// ModulePlanEntry is a lightweight stub used in topic.json's module plan.
// It records the planned modules without lesson details.
type ModulePlanEntry struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

// ModuleFile represents the module.json written/updated during Passes 2-3.
// Pass 2 populates metadata and learning objectives; Pass 3 adds the assessment.
type ModuleFile struct {
	ID                 string          `json:"id"`
	Title              string          `json:"title"`
	Description        string          `json:"description"`
	Order              int             `json:"order"`
	LearningObjectives []string        `json:"learning_objectives"`
	EstimatedMinutes   int             `json:"estimated_minutes"`
	Assessment         json.RawMessage `json:"assessment"`
}
