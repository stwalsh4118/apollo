package respond

import (
	"encoding/json"
	"net/http"
)

const contentTypeJSON = "application/json"

// JSON writes a JSON response with the given status code and data.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Error writes a JSON error response with the given status code and message.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, errorBody{Error: message})
}

type errorBody struct {
	Error string `json:"error"`
}
