package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func GetGroupRequests(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID string) {
	userID := r.Context().Value("userID").(string)
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

	err = CheckExistenceAndPermission(db, w, userID, groupID)
	if err != nil {
		return
	}

	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")

	users, err := sqlite.SelectGroupRequests(db, groupID, last, lastID)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: users,
	})
}

func CheckExistenceAndPermission(db *sql.DB, w http.ResponseWriter, userID, groupID string) error {
	var creatorID string
	err := db.QueryRow(`SELECT creator_id FROM groups WHERE id = ?`, groupID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		help.RespondNotFound(w, "This group doesn't exist")
		return err
	}
	if err != nil {
		help.RespondServerError(w)
		return err
	}

	if userID != creatorID {
		help.Respond(w, &models.Response{
			Code:    http.StatusForbidden,
			Message: "not creator",
		})

		return errors.New("forbidden")
	}

	return nil
}
