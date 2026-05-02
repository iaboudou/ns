package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"rtf/help"
	"rtf/models"
)

func DeleteGroup(w http.ResponseWriter, db *sql.DB, groupID, userID string) {
	var creatorID string
	err := db.QueryRow(`SELECT creator_id FROM groups WHERE id = ?`, groupID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		help.RespondNotFound(w, "This group doesn't exist")
		return
	}

	if creatorID != userID {
		help.Respond(w, &models.Response{
			Code:    http.StatusForbidden,
			Message: "Only the creator can delete the group",
		})
		return
	}

	if err != nil {
		help.RespondServerError(w)
		return
	}

	_, err = db.Exec(`DELETE FROM groups
			 WHERE id = ? AND creator_id = ?`, groupID, userID)
	if err != nil {
		fmt.Println("error while deleting group :", err)
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{Code: http.StatusOK})
}
