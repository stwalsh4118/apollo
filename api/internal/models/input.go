package models

import "encoding/json"

// TopicInput is the request body for creating or updating a topic.
type TopicInput struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Difficulty     string   `json:"difficulty,omitempty"`
	EstimatedHours float64  `json:"estimated_hours,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Status         string   `json:"status"`
	Version        int      `json:"version,omitempty"`
	SourceURLs     []string `json:"source_urls,omitempty"`
	GeneratedAt    string   `json:"generated_at,omitempty"`
	GeneratedBy    string   `json:"generated_by,omitempty"`
	ParentTopicID  string   `json:"parent_topic_id,omitempty"`
}

// ModuleInput is the request body for creating a module.
type ModuleInput struct {
	ID                 string          `json:"id"`
	TopicID            string          `json:"topic_id"`
	Title              string          `json:"title"`
	Description        string          `json:"description,omitempty"`
	LearningObjectives []string        `json:"learning_objectives,omitempty"`
	EstimatedMinutes   int             `json:"estimated_minutes,omitempty"`
	SortOrder          int             `json:"sort_order"`
	Assessment         json.RawMessage `json:"assessment,omitempty"`
}

// LessonInput is the request body for creating a lesson.
type LessonInput struct {
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

// ConceptInput is the request body for creating a concept.
type ConceptInput struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Definition      string   `json:"definition"`
	DefinedInLesson string   `json:"defined_in_lesson,omitempty"`
	DefinedInTopic  string   `json:"defined_in_topic,omitempty"`
	Difficulty      string   `json:"difficulty,omitempty"`
	FlashcardFront  string   `json:"flashcard_front,omitempty"`
	FlashcardBack   string   `json:"flashcard_back,omitempty"`
	Status          string   `json:"status,omitempty"`
	Aliases         []string `json:"aliases,omitempty"`
}

// ConceptReferenceInput is the request body for adding a concept reference.
// ConceptID is supplied via the URL path parameter, not the request body.
type ConceptReferenceInput struct {
	LessonID string `json:"lesson_id"`
	Context  string `json:"context,omitempty"`
}

// PrerequisiteInput is the request body for adding a topic prerequisite.
type PrerequisiteInput struct {
	TopicID             string `json:"topic_id"`
	PrerequisiteTopicID string `json:"prerequisite_topic_id"`
	Priority            string `json:"priority"`
	Reason              string `json:"reason,omitempty"`
}

// RelationInput is the request body for adding a topic relation.
type RelationInput struct {
	TopicA       string `json:"topic_a"`
	TopicB       string `json:"topic_b"`
	RelationType string `json:"relation_type"`
	Description  string `json:"description,omitempty"`
}
