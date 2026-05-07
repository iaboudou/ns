package handlers

import (
	"database/sql"
	"net/http"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func JoinGroup(w http.ResponseWriter, r *http.Request, hub *models.Hub, db *sql.DB, groupID, userID, Type string) {
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
		help.RespondServerError(w)
		return
	}

	CurrentUser := r.Context().Value("userID").(string)

	creator := models.User{}

	db.QueryRow(`
	SELECT u.id, u.nickname, u.firstname, u.lastname, u.profile_image 
	FROM users u
	JOIN group_members gm ON gm.user_id = u.id
	WHERE gm.group_id = ? 
	AND role = 'creator'
	`, groupID).
		Scan(&creator.ID, &creator.Nickname, &creator.Firstname, &creator.Lastname, &creator.ProfileImage)

	if CurrentUser == userID {
		requester, err := sqlite.GetUserByID(db, CurrentUser)
		if err != nil {
			help.RespondServerError(w)
			return
		}

		hub.Notify <- models.FollowNotif{
			FromUser:  requester,
			ToUserID:  creator.ID,
			NotifType: "group_request",
			GroupID:   groupID,
		}
	} else {
		hub.Notify <- models.FollowNotif{
			FromUser:  creator,
			ToUserID:  userID,
			NotifType: "group_invite",
		}
	}

	help.Respond(w, &models.Response{Code: 200})
}
