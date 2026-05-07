package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"rtf/help"
	"rtf/models"
)

func VoteEvent(w http.ResponseWriter, r *http.Request, db *sql.DB, eventID, userID string) {
	var exists int
	err := db.QueryRow(`SELECT 1 FROM events WHERE id = ?`, eventID).Scan(&exists)
	if err == sql.ErrNoRows {
		help.RespondNotFound(w, "This event doesn't exist")
		return
	}

	if err != nil {
		help.RespondServerError(w)
		return
	}

	var data map[string]string
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid body",
		})
		return
	}

	vote := strings.TrimSpace(data["vote"])

	if vote != "going" && vote != "notgoing" {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid vote",
		})
		return
	}

	_, err = db.Exec(`
			INSERT INTO event_responses (event_id, user_id, status)
			VALUES (?, ?, ?)
			ON CONFLICT(event_id, user_id)
			DO UPDATE SET status = excluded.status;`, eventID, userID, vote)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{Code: http.StatusOK})
}
