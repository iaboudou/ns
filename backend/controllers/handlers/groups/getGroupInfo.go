package handlers

import (
	"database/sql"
	"net/http"

	"rtf/help"
	"rtf/models"
)

func GetGroupBasicInfo(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID string) {
	user_id := r.Context().Value("userID").(string)

	info := models.GroupeInfo{IsCreator: false}

	var creator_id string
	err := db.QueryRow(`
	SELECT title, description, creator_id, image
	FROM groups
	WHERE id = ?
	`, groupID).Scan(&info.Title, &info.Description, &creator_id, &info.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			help.RespondNotFound(w, "This group doesn't exist")
			return
		}

		help.RespondServerError(w)
		return
	}

	if creator_id == user_id {
		info.IsCreator = true
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM group_members
	WHERE status ='accepted' AND group_id = ?`, groupID).Scan(&info.MemberAmount)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	// possibility to take unread message count and request number

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: info,
	})
}
