package controllers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"rtf/models"

	"github.com/gofrs/uuid/v5"
)

func SendFollowNotification(clients map[string][]*models.Client, db *sql.DB, fromUser models.User, toUserID, notifType, groupID string) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	if strings.Contains(notifType, "follow") {
		// if it already exist remove it to rewrite it
		db.Exec(`DELETE FROM notifications WHERE user_id = ? AND type = ? AND ref_id = ?`, toUserID, notifType, fromUser.ID)
	}

	_, err = db.Exec(`
		INSERT INTO notifications (id, user_id, type, ref_id, is_read, created_at, group_id)
		VALUES (?, ?, ?, ?, 0, ?, ?)
	`, id.String(), toUserID, notifType, fromUser.ID, time.Now().Unix(), groupID)
	if err != nil {
		return
	}

	// Push WS event if target user is online
	cs, online := clients[toUserID]
	if !online {
		return
	}

	payload := map[string]any{
		"event":      "notification",
		"id":         id.String(),
		"type":       notifType,
		"ref_id":     fromUser.ID,
		"from_name":  fromUser.Firstname + " " + fromUser.Lastname,
		"from_image": fromUser.ProfileImage,
		"created_at": time.Now().Unix(),
		"is_read":    false,
		"group_id":   groupID,
	}

	for _, c := range cs {
		c.Mu.Lock()
		c.Ws.WriteJSON(payload)
		c.Mu.Unlock()
	}
}

func GetNotificationsWS(clients map[string][]*models.Client, db *sql.DB, msg models.Message) error {

	userID := msg.SenderID

	rows, err := db.Query(`
		SELECT n.id, n.type, n.ref_id, n.is_read, n.created_at, n.group_id,
		       u.firstname, u.lastname, u.profile_image
		FROM notifications n
		LEFT JOIN users u ON u.id = n.ref_id
		WHERE n.user_id = ?
		ORDER BY n.created_at DESC
		LIMIT 50
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type NotifRow struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		RefID     string `json:"ref_id"`
		IsRead    bool   `json:"is_read"`
		CreatedAt int64  `json:"created_at"`
		FromName  string `json:"from_name"`
		FromImage string `json:"from_image"`
		GroupId   string `json:"group_id"`
	}

	notifs := []NotifRow{}
	for rows.Next() {
		var n NotifRow
		var isReadBool bool
		var refID, firstname, lastname, profileImage sql.NullString
		var groupID string
		var createdAtRaw interface{}

		err := rows.Scan(&n.ID, &n.Type, &refID, &isReadBool, &createdAtRaw, &firstname, &lastname, &profileImage, &groupID)
		if err != nil {
			fmt.Println("Scan error:", err)
			continue
		}

		n.RefID = refID.String

		n.IsRead = isReadBool
		n.FromName = firstname.String + " " + lastname.String
		n.FromImage = profileImage.String

		switch v := createdAtRaw.(type) {
		case int64:
			n.CreatedAt = v
		case time.Time:
			n.CreatedAt = v.Unix()
		case string:
			fmt.Sscanf(v, "%d", &n.CreatedAt)
		case []byte:
			fmt.Sscanf(string(v), "%d", &n.CreatedAt)
		}

		notifs = append(notifs, n)
	}

	db.Exec(`UPDATE notifications SET is_read = 1 WHERE user_id = ?`, userID)
	if cs, ok := clients[userID]; ok {
		for _, c := range cs {
			c.Mu.Lock()
			c.Ws.WriteJSON(map[string]any{
				"event": "notifications_list",
				"data":  notifs,
			})
			c.Mu.Unlock()
		}
	}
	return nil
}

func GetUnreadNotificationCountWS(clients map[string][]*models.Client, db *sql.DB, msg models.Message) error {
	userID := msg.SenderID

	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0`, userID).Scan(&count)
	if err != nil {
		return err
	}

	if cs, ok := clients[userID]; ok {
		for _, c := range cs {
			c.Mu.Lock()
			c.Ws.WriteJSON(map[string]any{
				"event": "unread_notifications_count",
				"count": count,
			})
			c.Mu.Unlock()
		}
	}
	return nil
}
