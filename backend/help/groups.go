package help

import (
	"net/http"
	"strings"

	"rtf/models"
)

func GetPathSegments(path string) []string {
	clean := strings.TrimPrefix(path, "/api/groups")
	if clean == "" {
		return []string{}
	}

	parts := strings.Split(clean, "/")
	var segments []string
	for _, p := range parts {
		if p != "" {
			segments = append(segments, p)
		}
	}
	return segments
}

// /
func RespondNotFound(w http.ResponseWriter, msg string) {
	Respond(w, &models.Response{
		Code:    http.StatusNotFound,
		Message: msg,
	})
}

// /
func RespondMethodNotAllowed(w http.ResponseWriter) {
	Respond(w, &models.Response{
		Code: http.StatusMethodNotAllowed,
	})
}
