package models

import "encoding/json"

// ModuleSummary is the brief representation used inside topic views.
type ModuleSummary struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description,omitempty"`
	EstimatedMinutes int    `json:"estimated_minutes,omitempty"`
	SortOrder        int    `json:"sort_order"`
}

// moduleBase holds fields shared between ModuleDetail and ModuleFull.
type moduleBase struct {
	ID                 string          `json:"id"`
	TopicID            string          `json:"topic_id"`
	Title              string          `json:"title"`
	Description        string          `json:"description,omitempty"`
	LearningObjectives []string        `json:"learning_objectives,omitempty"`
	EstimatedMinutes   int             `json:"estimated_minutes,omitempty"`
	SortOrder          int             `json:"sort_order"`
	Assessment         json.RawMessage `json:"assessment,omitempty"`
}

// ModuleDetail includes the module's lessons.
type ModuleDetail struct {
	moduleBase
	Lessons []LessonSummary `json:"lessons"`
}

// ModuleFull includes lessons with full content (for topic full tree).
type ModuleFull struct {
	moduleBase
	Lessons []LessonFull `json:"lessons"`
}
