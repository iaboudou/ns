package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"rtf/help"
	"rtf/models"
)

func InviteUser(w http.ResponseWriter, r *http.Request, hub *models.Hub, db *sql.DB, groupID string) {
	body := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid body",
		})
		return
	}

	invitedUserID := body["invitedUser"]

	var exists int
	err = db.QueryRow(`SELECT 1 FROM users WHERE id = ?`, invitedUserID).Scan(&exists)
	if err == sql.ErrNoRows {
		help.RespondNotFound(w, "This user doesn't exist")
		return
	}

	if err != nil {
		help.RespondServerError(w)
		return
	}

	JoinGroup(w, r, hub, db, groupID, invitedUserID, "invite")
}
