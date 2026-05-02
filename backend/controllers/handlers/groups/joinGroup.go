package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"rtf/help"
	"rtf/models"
)

func JoinGroup(w http.ResponseWriter, db *sql.DB, groupID, userID, Type string) {
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

	_, err = db.Exec(`INSERT INTO group_members (group_id, user_id, type) VALUES (?, ?, ?)`, groupID, userID, Type)
	if err != nil {
		fmt.Println("error while adding group request:", err)
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{Code: 200})
}
