package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"rtf/help"
	"rtf/models"
)

func LeaveGroup(w http.ResponseWriter, db *sql.DB, groupID, userID string) {
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

	_, err = db.Exec(`DELETE FROM group_members
	WHERE group_id = ? AND user_id = ?`, groupID, userID)
	if err != nil {
		fmt.Println("error while removing a user from a group: ", err)
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{Code: http.StatusOK})
}
