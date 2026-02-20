package models

import "encoding/json"

// ResearchJobStatus is the status of a research job.
type ResearchJobStatus = string

// Research job status constants matching the DB CHECK constraint.
const (
	ResearchStatusQueued      ResearchJobStatus = "queued"
	ResearchStatusResearching ResearchJobStatus = "researching"
	ResearchStatusResolving   ResearchJobStatus = "resolving"
	ResearchStatusPublished   ResearchJobStatus = "published"
	ResearchStatusFailed      ResearchJobStatus = "failed"
	ResearchStatusCancelled   ResearchJobStatus = "cancelled"
)

// ResearchJob represents a row in the research_jobs table.
type ResearchJob struct {
	ID           string          `json:"id"`
	RootTopic    string          `json:"root_topic"`
	CurrentTopic string          `json:"current_topic,omitempty"`
	Status       string          `json:"status"`
	Progress     json.RawMessage `json:"progress,omitempty"`
	Error        string          `json:"error,omitempty"`
	StartedAt    string          `json:"started_at,omitempty"`
	CompletedAt  string          `json:"completed_at,omitempty"`
}

// ResearchProgress tracks the current state of a research pipeline execution.
type ResearchProgress struct {
	CurrentPass      int            `json:"current_pass"`
	TotalPasses      int            `json:"total_passes"`
	ModulesPlanned   int            `json:"modules_planned"`
	ModulesCompleted int            `json:"modules_completed"`
	ConceptsFound    int            `json:"concepts_found"`
	PassDescriptions map[int]string `json:"pass_descriptions,omitempty"`
}

// CreateResearchJobInput is the request body for POST /api/research.
// Topic is required and must be non-empty; validated at the handler layer.
type CreateResearchJobInput struct {
	Topic string `json:"topic"`
	Brief string `json:"brief,omitempty"`
}

// ResearchJobSummary is a subset of ResearchJob for list responses.
type ResearchJobSummary struct {
	ID          string `json:"id"`
	RootTopic   string `json:"root_topic"`
	Status      string `json:"status"`
	StartedAt   string `json:"started_at,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// IsTerminalStatus returns true if the research job status is a terminal state.
func IsTerminalStatus(status string) bool {
	return status == ResearchStatusPublished ||
		status == ResearchStatusFailed ||
		status == ResearchStatusCancelled
}
