package handlers

import (
	"database/sql"
	"net/http"

	"rtf/help"
	"rtf/models"
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

	CurrentUserID := r.Context().Value("userID").(string)

	var groupCreatorID string

	err = db.QueryRow(`
	SELECT u.id
	FROM users u
	JOIN group_members gm ON gm.user_id = u.id
	WHERE gm.group_id = ? 
	AND role = 'creator'
	`, groupID).Scan(&groupCreatorID)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	if CurrentUserID == groupCreatorID {
		hub.Notif <- models.Notification{
			SenderID:   groupCreatorID,
			ReceiverID: userID,
			Type:       "group_invite",
			GroupID:    groupID,
		}
	} else {
		hub.Notif <- models.Notification{
			SenderID:   userID,
			ReceiverID: groupCreatorID,
			Type:       "group_request",
			GroupID:    groupID,
		}
	}

	help.Respond(w, &models.Response{Code: 200})
}
