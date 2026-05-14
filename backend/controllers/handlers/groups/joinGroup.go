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

	currentUserID := r.Context().Value("userID").(string)
	err = sqlite.InsertGroupMember(db, userID, groupID, currentUserID, Type)
	if err != nil {
		if err.Error() == "already member" {
			help.Respond(w, &models.Response{
				Code:    http.StatusConflict,
				Message: err.Error(),
			})
			return
		}

		help.RespondServerError(w)
		return
	}

	err = SendNotif(hub, db, groupID, userID, currentUserID, Type)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{Code: 200})
}

func SendNotif(hub *models.Hub, db *sql.DB, groupID, userID, currentUserID, Type string) error {
	var groupCreatorID string

	err := db.QueryRow(`
	SELECT u.id
	FROM users u
	JOIN group_members gm ON gm.user_id = u.id
	WHERE gm.group_id = ? 
	AND role = 'creator'
	`, groupID).Scan(&groupCreatorID)
	if err != nil {
		return err
	}

	if Type == "invite" {
		hub.Notif <- models.Notification{
			SenderID:   currentUserID,
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

	return nil
}
