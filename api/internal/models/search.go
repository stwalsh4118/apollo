package models

// SearchResult represents a single full-text search hit.
type SearchResult struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Title      string `json:"title"`
	Snippet    string `json:"snippet"`
}
