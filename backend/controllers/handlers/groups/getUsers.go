package handlers

import (
	"database/sql"
	"net/http"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID, userId string) {
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

	var isMember bool
	err = db.QueryRow(`SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?`, userId, groupID).Scan(&isMember)
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
	search := r.URL.Query().Get("search")

	users, err := sqlite.SelectFollowers(db, userId, groupID, last, lastID, search)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: 200,
		Data: users,
	})
}
