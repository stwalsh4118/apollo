package models

import "encoding/json"

// ParseJSONStringSlice parses a JSON TEXT column into a Go string slice.
// Returns nil for NULL or empty values.
func ParseJSONStringSlice(raw *string) []string {
	if raw == nil || *raw == "" || *raw == "null" {
		return nil
	}

	var result []string
	if err := json.Unmarshal([]byte(*raw), &result); err != nil {
		return nil
	}

	return result
}

// ParseJSONRaw converts a nullable JSON TEXT column into json.RawMessage.
// Returns nil for NULL or empty values.
func ParseJSONRaw(raw *string) json.RawMessage {
	if raw == nil || *raw == "" || *raw == "null" {
		return nil
	}

	return json.RawMessage(*raw)
}
