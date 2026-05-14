package handlers

import (
	"database/sql"
	"net/http"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func GetGroupBasicInfo(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID string) {
	user_id := r.Context().Value("userID").(string)
	var isMember bool
	err := db.QueryRow(`SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?`, user_id, groupID).Scan(&isMember)
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

	info := models.GroupeInfo{IsCreator: false}

	err = sqlite.SelectGroupDescription(db, &info, groupID, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			help.RespondNotFound(w, "This group doesn't exist")
			return
		}

		if err.Error() == "not a member" {
			help.Respond(w, &models.Response{
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			return
		}

		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: info,
	})
}
