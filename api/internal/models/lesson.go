package models

import "encoding/json"

// LessonSummary is the brief representation used inside module views.
type LessonSummary struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	SortOrder        int    `json:"sort_order"`
	EstimatedMinutes int    `json:"estimated_minutes,omitempty"`
}

// lessonBase holds fields shared between LessonDetail and LessonFull.
type lessonBase struct {
	ID               string          `json:"id"`
	ModuleID         string          `json:"module_id"`
	Title            string          `json:"title"`
	SortOrder        int             `json:"sort_order"`
	EstimatedMinutes int             `json:"estimated_minutes,omitempty"`
	Content          json.RawMessage `json:"content"`
	Examples         json.RawMessage `json:"examples,omitempty"`
	Exercises        json.RawMessage `json:"exercises,omitempty"`
	ReviewQuestions  json.RawMessage `json:"review_questions,omitempty"`
}

// LessonDetail is the full lesson with all content fields.
type LessonDetail struct {
	lessonBase
}

// LessonFull is used in the topic full tree, includes concepts.
type LessonFull struct {
	lessonBase
	Concepts []ConceptSummary `json:"concepts,omitempty"`
}
