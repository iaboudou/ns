package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"rtf/help"
	"rtf/models"
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

	var request struct {
		Sender   string `json:"sender"`
		Decicion string `json:"decision"`
	}

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

	switch request.Decicion {
	case "rejected":
		_, err = db.Exec(`DELETE FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, request.Sender)
	case "accepted":
		_, err = db.Exec(`UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ?`, request.Decicion, groupID, request.Sender)
	default:
		help.Respond(w, &models.Response{Code: http.StatusBadRequest, Message: "unknown decision"})
		return
	}

	if err != nil {
		fmt.Println("error while managing group request:", err)
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{Code: http.StatusOK})
}
