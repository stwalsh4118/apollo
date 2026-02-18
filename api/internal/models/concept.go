package models

// ConceptSummary is the brief representation used in list views.
type ConceptSummary struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Definition     string   `json:"definition"`
	Difficulty     string   `json:"difficulty,omitempty"`
	Status         string   `json:"status"`
	DefinedInTopic string   `json:"defined_in_topic,omitempty"`
	Aliases        []string `json:"aliases,omitempty"`
}

// ConceptDetail includes references and full fields.
type ConceptDetail struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Definition      string             `json:"definition"`
	DefinedInLesson string             `json:"defined_in_lesson,omitempty"`
	DefinedInTopic  string             `json:"defined_in_topic,omitempty"`
	Difficulty      string             `json:"difficulty,omitempty"`
	FlashcardFront  string             `json:"flashcard_front,omitempty"`
	FlashcardBack   string             `json:"flashcard_back,omitempty"`
	Status          string             `json:"status"`
	Aliases         []string           `json:"aliases,omitempty"`
	References      []ConceptReference `json:"references,omitempty"`
}

// ConceptReference links a concept to a lesson where it's referenced.
type ConceptReference struct {
	LessonID    string `json:"lesson_id"`
	LessonTitle string `json:"lesson_title"`
	Context     string `json:"context,omitempty"`
}
