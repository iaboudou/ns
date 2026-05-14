package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func ManageRequest(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID, userID string) {
	var exists int
	err := db.QueryRow(`SELECT 1 FROM groups WHERE id = ?`, groupID).Scan(&exists)
	if err == sql.ErrNoRows {
		help.RespondNotFound(w, "This group doesn't exist")
		return
	}
	if err != nil {
		help.RespondServerError(w)
		return
	}

	request := models.Request{}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid body",
		})
		return
	}

	if request.Sender == "" {
		request.Sender = userID
	} else {
		err = db.QueryRow(`SELECT 1 FROM users WHERE id = ?`, request.Sender).Scan(&exists)
		if err == sql.ErrNoRows {
			help.RespondNotFound(w, "This user doesn't exist")
			return
		}
		if err != nil {
			help.RespondServerError(w)
			return
		}
	}

	err = sqlite.UpdateRequest(db, &request, groupID)
	if err != nil {
		switch err.Error() {
		case "unknown decision":
			help.Respond(w, &models.Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		case "":
			help.RespondNotFound(w, "ressource doesn't exists")
		default:
			help.RespondServerError(w)
		}
		return
	}

	err = sqlite.RemoveNotif(db, &request, userID, groupID)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{Code: http.StatusOK})
}
