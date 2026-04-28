package pkg

import (
	"encoding/json"
	"net/http"
)

func RespondOK(w http.ResponseWriter, rep any, TYPE string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if rep != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			TYPE:      rep,
			"success": "done successfully",
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": "done successfully",
		})
	}
}

// RespondNotOK returns an HTTP error response based on the provided error type
func RespondNotOK(w http.ResponseWriter, Type string) {
	switch Type {
	case "unauthorized":
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	case "forbidden":
		http.Error(w, "action not allowed", http.StatusForbidden)
	case "server-error":
		http.Error(w, "server error", http.StatusInternalServerError)
	case "notallowed":
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	case "badrequest":
		http.Error(w, "action not allowed", http.StatusBadRequest)
	case "notfound":
		http.Error(w, "not found", http.StatusNotFound)
	}
}
