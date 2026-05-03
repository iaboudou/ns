package sqlite

import (
	"database/sql"
	"time"

	"rtf/models"

	"github.com/gofrs/uuid/v5"
)

// func SelectOtherUsers(db *sql.DB, clients map[string][]*models.Client, currentUserID string) ([]models.OtherClient, error) {
// 	rows, err := db.Query(`SELECT nickname, id FROM user WHERE id != ?`, currentUserID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	users := []models.OtherClient{}

// 	for rows.Next() {
// 		var u models.OtherClient
// 		var id string
// 		if err := rows.Scan(&u.NickName, &id); err != nil {
// 			return nil, err
// 		}

// 		err := db.QueryRow(`
// 				SELECT created_at
// 				FROM private_message
// 				WHERE (sender_id = ? AND receiver_id = ?)
// 				OR (receiver_id = ? AND sender_id = ?)
// 				ORDER BY created_at DESC
// 				LIMIT 1
// 				`, currentUserID, id, currentUserID, id).Scan(&u.LastChat)
// 		if err != nil && err != sql.ErrNoRows {
// 			return nil, err
// 		}

// 		err = db.QueryRow(`
//     			SELECT COUNT(*)
//     			FROM private_message
//     			WHERE sender_id = ?
//       			AND receiver_id = ?
//      			AND is_read = FALSE
// 				`, id, currentUserID).Scan(&u.Pending_Message)
// 		if err != nil && err != sql.ErrNoRows {
// 			return nil, err
// 		}

// 		cs, exists := clients[u.NickName]
// 		u.Online = exists && len(cs) > 0
// 		users = append(users, u)
// 	}

// 	return users, nil
// }

func MarkRead(db *sql.DB, receiverID, senderID string) error {
	_, err := db.Exec(`
		UPDATE messages
		SET is_not_read = FALSE
		WHERE sender_id = ?
		AND receiver_id = ?
	`, senderID, receiverID)
	return err
}

func SelectUnreadCount(db *sql.DB, msg *models.Message) (int, error) {
	var amount int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM private_message pm
		JOIN user s ON s.id = pm.sender_id
		JOIN user r ON r.id = pm.receiver_id
		WHERE s.nickname = ?
		AND r.nickname = ?
		AND pm.is_read = FALSE
	`, msg.ReceiverID, msg.SenderID).Scan(&amount)

	return amount, err
}

func SelectOldMessages(db *sql.DB, msg *models.Message) ([]map[string]any, error) {
	if msg.LastReadTime == 0 {
		msg.LastReadTime = time.Now().UnixMilli()
	}

	rows, err := db.Query(`
		SELECT 
			m.id, m.created_at, m.content, m.sender_id, m.receiver_id,
			s.nickname, s.firstname, s.lastname, s.profile_image,
			r.nickname, r.firstname, r.lastname, r.profile_image
		FROM messages m
		JOIN users s ON s.id = m.sender_id
		JOIN users r ON r.id = m.receiver_id
		WHERE ((m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?))
		  AND m.created_at < ?
		ORDER BY m.created_at DESC
		LIMIT 10
	`, msg.SenderID, msg.ReceiverID, msg.ReceiverID, msg.SenderID, msg.LastReadTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []map[string]any{}
	for rows.Next() {
		var ID, content, senderID, receiverID string
		var senderNickname, senderFirstname, senderLastname, senderProfile string
		var receiverNickname, receiverFirstname, receiverLastname, receiverProfile string
		var createdAt int64

		if err := rows.Scan(
			&ID, &createdAt, &content, &senderID, &receiverID,
			&senderNickname, &senderFirstname, &senderLastname, &senderProfile,
			&receiverNickname, &receiverFirstname, &receiverLastname, &receiverProfile,
		); err != nil {
			continue
		}

		messages = append(messages, map[string]any{
			"id":                ID,
			"content":           content,
			"created_at":        createdAt,
			"sender_id":         senderID,
			"receiver_id":       receiverID,
			"sender_nickname":   senderNickname,
			"sender_fullname":   senderFirstname + " " + senderLastname,
			"sender_profile":    senderProfile,
			"receiver_nickname": receiverNickname,
			"receiver_fullname": receiverFirstname + " " + receiverLastname,
			"receiver_profile":  receiverProfile,
		})
	}

	return messages, nil
}

func InsertNewMessage(db *sql.DB, msg *models.Message) error {
	messageID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	msg.ID = messageID.String()

	_, err = db.Exec(`
		INSERT INTO messages (id, sender_id, receiver_id, content, created_at)
		VALUES ( ?, ?, ?, ?, ?)
	`, messageID.String(), msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt)

	return err
}

func GetUserByID(db *sql.DB, userID string) (models.User, error) {
	u := models.User{ID: userID}

	err := db.QueryRow(`SELECT nickname, firstname, lastname, profile_image FROM users
	WHERE id = ?`, userID).Scan(&u.Nickname, &u.Firstname, &u.Lastname, &u.ProfileImage)
	if err != nil {
		return models.User{}, err
	}

	return u, err
}
