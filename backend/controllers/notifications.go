package controllers

import (
	"database/sql"
	"time"

	"rtf/models"
	"rtf/pkg/db/sqlite"

	"github.com/gofrs/uuid/v5"
)

func GetNotificationsWS(clients map[string][]*models.Client, db *sql.DB, userID string) error {
	rows, err := db.Query(`
		SELECT n.id, n.type, n.group_id, n.created_at, nu.is_read,
		       u.id, u.nickname, u.firstname, u.lastname, u.profile_image
		FROM notifications n
		JOIN notification_users nu ON nu.notification_id = n.id
		JOIN users u ON u.id = n.sender_id
		WHERE nu.user_id = ?
		ORDER BY n.created_at DESC
	`, userID)
	if err != nil {
		return err
	}

	defer rows.Close()

	notifications := []map[string]any{}

	for rows.Next() {
		var id, notifType, groupID string
		var createdAt time.Time
		var isRead int
		var senderID, senderNickname, senderFirstname, senderLastname, senderProfile string

		err := rows.Scan(&id, &notifType, &groupID, &createdAt, &isRead, &senderID, &senderNickname, &senderFirstname, &senderLastname, &senderProfile)
		if err != nil {
			return err
		}

		notifications = append(notifications, map[string]any{
			"id":              id,
			"type":            notifType,
			"group_id":        groupID,
			"created_at":      createdAt,
			"is_read":         isRead == 1,
			"sender_id":       senderID,
			"sender_nickname": senderNickname,
			"sender_fullname": senderFirstname + " " + senderLastname,
			"sender_profile":  senderProfile,
		})
	}

	if cs, ok := clients[userID]; ok {
		for _, c := range cs {
			c.Mu.Lock()
			c.Ws.WriteJSON(map[string]any{
				"event":         "notifications",
				"notifications": notifications,
			})
			c.Mu.Unlock()
		}
	}
	return nil
}

func BroadCastEventCreation(db *sql.DB, clients map[string][]*models.Client, notif *models.Notification) error {
	notifID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format("2006-01-02 15:04:05")

	rows, err := db.Query(`
		SELECT user_id 
		FROM group_members
		WHERE group_id = ? AND status = 'accepted'
	`, notif.GroupID)
	if err != nil {
		return err
	}

	members := map[string]bool{}
	var receivers []string

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return err
		}
		members[uid] = true
		if uid != notif.SenderID {
			receivers = append(receivers, uid)
		}
	}
	rows.Close()

	if err := sqlite.InsertNotifInDB(db, notif, notifID, now, receivers); err != nil {
		return err
	}

	sender, err := sqlite.GetUserByID(db, notif.SenderID)
	if err != nil {
		return err
	}

	for clientID, clientArr := range clients {
		if clientID == notif.SenderID || !members[clientID] {
			continue
		}

		for _, c := range clientArr {
			c.Mu.Lock()
			err := c.Ws.WriteJSON(map[string]any{
				"event": "new_notification",
				"notif": map[string]any{
					"id":              notifID,
					"type":            "event",
					"is_read":         false,
					"sender_id":       sender.ID,
					"sender_nickname": sender.Nickname,
					"sender_fullname": sender.Firstname + " " + sender.Lastname,
					"sender_profile":  sender.ProfileImage,
					"created_at":      now,
					"group_id":        notif.GroupID,
				},
			})
			c.Mu.Unlock()

			if err != nil {
				c.Ws.Close()
			}
		}
	}

	return nil
}

func Notify(db *sql.DB, clients map[string][]*models.Client, notif *models.Notification) error {
	// check if the same notification already exists to avoid spam notif
	var exists int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM notifications n
		JOIN notification_users nu ON nu.notification_id = n.id
		WHERE n.type = ?
		AND n.sender_id = ?
		AND nu.user_id = ?
		AND n.group_id IS ?
	`, notif.Type, notif.SenderID, notif.ReceiverID, notif.GroupID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}

	notifID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format("2006-01-02 15:04:05")

	if err := sqlite.InsertNotifInDB(db, notif, notifID, now, []string{notif.ReceiverID}); err != nil {
		return err
	}

	sender, err := sqlite.GetUserByID(db, notif.SenderID)
	if err != nil {
		return err
	}

	if cs, ok := clients[notif.ReceiverID]; ok {
		for _, c := range cs {
			c.Mu.Lock()
			err := c.Ws.WriteJSON(map[string]any{
				"event": "new_notification",
				"notif": map[string]any{
					"id":              notifID,
					"type":            notif.Type,
					"is_read":         false,
					"sender_id":       sender.ID,
					"sender_nickname": sender.Nickname,
					"sender_fullname": sender.Firstname + " " + sender.Lastname,
					"sender_profile":  sender.ProfileImage,
					"receiver_id":     notif.ReceiverID,
					"created_at":      now,
					"group_id":        notif.GroupID,
				},
			})
			c.Mu.Unlock()

			if err != nil {
				c.Ws.Close()
			}
		}
	}
	return nil
}
