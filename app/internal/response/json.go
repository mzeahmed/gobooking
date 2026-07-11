// Package response provides helper functions to write HTTP responses.
//
// Keeping all response helpers in a dedicated package ensures
// consistency across the application and avoids duplicating
// JSON encoding logic inside individual handlers.
package response

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response.
//
// The helper automatically sets the Content-Type header,
// writes the HTTP status code and serializes the provided
// value using the standard JSON encoder.
func JSON(w http.ResponseWriter, status int, data any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}
