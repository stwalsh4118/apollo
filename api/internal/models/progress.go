package models

// Progress status constants matching the DB CHECK constraint.
const (
	ProgressStatusNotStarted = "not_started"
	ProgressStatusInProgress = "in_progress"
	ProgressStatusCompleted  = "completed"
)

// LessonProgress represents a single lesson's learning progress.
type LessonProgress struct {
	LessonID    string `json:"lesson_id"`
	LessonTitle string `json:"lesson_title,omitempty"`
	Status      string `json:"status"`
	StartedAt   string `json:"started_at,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

// TopicProgress is the response for GET /api/progress/topics/:id.
type TopicProgress struct {
	TopicID string           `json:"topic_id"`
	Lessons []LessonProgress `json:"lessons"`
}

// ProgressSummary is the response for GET /api/progress/summary.
type ProgressSummary struct {
	TotalLessons         int     `json:"total_lessons"`
	CompletedLessons     int     `json:"completed_lessons"`
	CompletionPercentage float64 `json:"completion_percentage"`
	ActiveTopics         int     `json:"active_topics"`
}

// UpdateProgressInput is the request body for PUT /api/progress/lessons/:id.
type UpdateProgressInput struct {
	Status string `json:"status"`
	Notes  string `json:"notes,omitempty"`
}
