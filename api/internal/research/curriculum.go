package research

import "encoding/json"

// CurriculumOutput is the top-level structure matching the curriculum JSON schema.
type CurriculumOutput struct {
	ID             string              `json:"id"`
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	Difficulty     string              `json:"difficulty"`
	EstimatedHours float64             `json:"estimated_hours"`
	Tags           []string            `json:"tags"`
	Prerequisites  PrerequisitesOutput `json:"prerequisites"`
	RelatedTopics  []string            `json:"related_topics"`
	Modules        []ModuleOutput      `json:"modules"`
	SourceURLs     []string            `json:"source_urls"`
	GeneratedAt    string              `json:"generated_at"`
	Version        int                 `json:"version"`
}

// PrerequisitesOutput holds the three priority levels of prerequisites.
type PrerequisitesOutput struct {
	Essential      []PrerequisiteItem `json:"essential"`
	Helpful        []PrerequisiteItem `json:"helpful"`
	DeepBackground []PrerequisiteItem `json:"deep_background"`
}

// PrerequisiteItem is a single prerequisite reference.
type PrerequisiteItem struct {
	TopicID string `json:"topic_id"`
	Reason  string `json:"reason"`
}

// ModuleOutput matches the Module definition in the schema.
type ModuleOutput struct {
	ID                 string          `json:"id"`
	Title              string          `json:"title"`
	Description        string          `json:"description"`
	LearningObjectives []string        `json:"learning_objectives"`
	EstimatedMinutes   int             `json:"estimated_minutes"`
	Order              int             `json:"order"`
	Lessons            []LessonOutput  `json:"lessons"`
	Assessment         json.RawMessage `json:"assessment"`
}

// LessonOutput matches the Lesson definition in the schema.
type LessonOutput struct {
	ID                 string             `json:"id"`
	Title              string             `json:"title"`
	Order              int                `json:"order"`
	EstimatedMinutes   int                `json:"estimated_minutes"`
	Content            json.RawMessage    `json:"content"`
	ConceptsTaught     []ConceptTaughtOut `json:"concepts_taught"`
	ConceptsReferenced []ConceptRefOut    `json:"concepts_referenced"`
	Examples           json.RawMessage    `json:"examples"`
	Exercises          json.RawMessage    `json:"exercises"`
	ReviewQuestions    json.RawMessage    `json:"review_questions"`
}

// ConceptTaughtOut represents a concept defined in this lesson.
type ConceptTaughtOut struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Definition string       `json:"definition"`
	Flashcard  FlashcardOut `json:"flashcard"`
}

// FlashcardOut holds the spaced repetition card for a concept.
type FlashcardOut struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

// ConceptRefOut is a reference to a concept defined elsewhere.
type ConceptRefOut struct {
	ID        string `json:"id"`
	DefinedIn string `json:"defined_in"`
}
