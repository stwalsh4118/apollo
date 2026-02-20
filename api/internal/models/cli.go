package models

import "encoding/json"

// CLIResponse represents the JSON output from `claude -p --output-format json`.
type CLIResponse struct {
	Type             string          `json:"type"`
	SessionID        string          `json:"session_id"`
	Result           string          `json:"result,omitempty"`
	StructuredOutput json.RawMessage `json:"structured_output,omitempty"`
	IsError          bool            `json:"is_error"`
	ErrorMessage     string          `json:"error_message,omitempty"`
}
