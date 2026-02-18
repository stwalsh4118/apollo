package models

// GraphNode represents a node in the knowledge graph.
type GraphNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"` // "topic" or "concept"
}

// GraphEdge represents a relationship between two nodes.
type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"` // prerequisite, reference, or relation type value
}

// GraphData contains the full graph response with nodes and edges.
type GraphData struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}
