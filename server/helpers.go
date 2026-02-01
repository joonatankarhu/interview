package server

import (
	"encoding/json"
	"net/http"
)

// WriteError sends a JSON error response with the given message and status code.
func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
