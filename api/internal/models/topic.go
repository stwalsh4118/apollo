package models

// TopicSummary is the list-view representation of a topic.
type TopicSummary struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Difficulty     string   `json:"difficulty,omitempty"`
	EstimatedHours float64  `json:"estimated_hours,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Status         string   `json:"status"`
	ModuleCount    int      `json:"module_count"`
}

// topicBase holds fields shared between TopicDetail and TopicFull.
type topicBase struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Difficulty     string   `json:"difficulty,omitempty"`
	EstimatedHours float64  `json:"estimated_hours,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Status         string   `json:"status"`
	Version        int      `json:"version"`
	SourceURLs     []string `json:"source_urls,omitempty"`
	GeneratedAt    string   `json:"generated_at,omitempty"`
	GeneratedBy    string   `json:"generated_by,omitempty"`
	ParentTopicID  string   `json:"parent_topic_id,omitempty"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// TopicDetail includes the topic's modules without lesson content.
type TopicDetail struct {
	topicBase
	Modules []ModuleSummary `json:"modules"`
}

// TopicFull is the full nested tree: topic → modules → lessons → concepts.
type TopicFull struct {
	topicBase
	Modules []ModuleFull `json:"modules"`
}
