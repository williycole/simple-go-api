package api

import (
	"encoding/json"
	"net/http"
)

// HelloResponse is used for successful JSON responses with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ReverseRequest represents a request payload for reversing text.
type ReverseRequest struct {
	Text string `json:"text"`
}

// ReverseResponse represents the response payload for a reversed text request.
type ReverseResponse struct {
	Reversed string `json:"reversed"`
}

// ErrorResponse is used for JSON error responses with an error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

// respondWithError sends a JSON error response with a given status code.
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// respondWithJSON sends a JSON response with a given status code and payload.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
