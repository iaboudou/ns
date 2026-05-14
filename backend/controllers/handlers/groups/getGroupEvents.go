package handlers

import (
	"database/sql"
	"net/http"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func GetGroupEvents(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID, userID string) {
	var isMember bool
	err := db.QueryRow(`SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?`, userID, groupID).Scan(&isMember)
	if err != nil {
		if err == sql.ErrNoRows {
			help.Respond(w, &models.Response{
				Code:    http.StatusBadRequest,
				Message: "you are not member of the group",
			})
			return
		}

		help.RespondServerError(w)
		return
	}

	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")

	events, err := sqlite.SelectEventsInDB(db, userID, groupID, last, lastID)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: events,
	})
}
