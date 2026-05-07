package sqlite

import (
	"database/sql"
	"time"

	"rtf/models"

	"github.com/gofrs/uuid/v5"
)

func MarkRead(db *sql.DB, receiverID, senderID string) error {
	_, err := db.Exec(`
		UPDATE messages
		SET is_not_read = 0
		WHERE sender_id = ?
		AND receiver_id = ?
	`, senderID, receiverID)
	return err
}

func SelectUnreadCount(db *sql.DB, msg *models.Message) (int, error) {
	var amount int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM messages pm
		WHERE pm.sender_id = ?
		AND pm.receiver_id = ?
		AND pm.is_not_read = 1
	`, msg.SenderID, msg.ReceiverID).Scan(&amount)

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
		INSERT INTO messages (id, sender_id, receiver_id, content, is_not_read, created_at)
		VALUES (?, ?, ?, ?, 1, ?)
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

func GetChatUsers(db *sql.DB, userID string) ([]map[string]any, error) {
	rows, err := db.Query(`
		SELECT u.id, u.firstname, u.lastname, u.profile_image,
               IFNULL(MAX(m.created_at), 0) as last_msg,
               (SELECT COUNT(*) FROM messages WHERE sender_id = u.id AND receiver_id = ? AND is_not_read = 1) as unread_count
		FROM users u
        JOIN (
            SELECT following_id as user_id FROM followers WHERE follower_id = ? AND status = 'accepted'
            UNION
            SELECT follower_id as user_id FROM followers WHERE following_id = ? AND status = 'accepted'
        ) rel ON u.id = rel.user_id
        LEFT JOIN messages m ON (
            (m.sender_id = ? AND m.receiver_id = u.id)
            OR
            (m.sender_id = u.id AND m.receiver_id = ?)
        )
		WHERE u.id != ?
        GROUP BY u.id
		ORDER BY last_msg DESC, u.firstname ASC, u.lastname ASC
	`, userID, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []map[string]any{}
	for rows.Next() {
		var id, firstname, lastname, profileImage string
		var lastMsg int64
		var unreadCount int
		if err := rows.Scan(&id, &firstname, &lastname, &profileImage, &lastMsg, &unreadCount); err != nil {
			return nil, err
		}
		users = append(users, map[string]any{
			"id":            id,
			"firstname":     firstname,
			"lastname":      lastname,
			"profile_image": profileImage,
			"last_msg":      lastMsg,
			"unread_count":  unreadCount,
		})
	}
	return users, rows.Err()
}

func (r *Repo) GetUserAuthernameByID(userID string) (string, error) {
	var nickname string
	er := r.Db.QueryRow("SELECT nickname FROM users WHERE id=?", userID).Scan(&nickname)
	return nickname, er
}

func (r *Repo) SetMessageRead(senderID, receiverID string) error {
	_, er := r.Db.Exec(` UPDATE messages SET is_NOT_read = 0 WHERE sender_id = ? AND receiver_id = ? AND is_NOT_read = 1`, senderID, receiverID)
	if er != nil {
		return er
	}
	return nil
}

// this method get the messages history between two users from DB in order from the offset provided at ascending order
func (r *Repo) GetMessagesHistoryBetweenTwoUsers(senderID, receiverID string, offset int) ([]models.Message, error) {
	rows, er := r.Db.Query(`
        SELECT sender_id, receiver_id, content, created_at 
        FROM messages
        WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at DESC
        LIMIT 10 OFFSET ?`,
		senderID, receiverID, receiverID, senderID, offset)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	msgs := []models.Message{}
	for rows.Next() {
		m := models.Message{}
		er := rows.Scan(&m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt)
		if er != nil {
			continue
		}
		m.ReceiverName, er = r.GetUserAuthernameByID(m.ReceiverID)
		if er != nil {
			continue
		}
		m.SenderName, er = r.GetUserAuthernameByID(m.SenderID)
		if er != nil {
			continue
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// GetChatUsersDB returns all users that have any follow relationship with userID
// (either userID follows them, or they follow userID — at least one side accepted)
func (r *Repo) GetChatUsersDB(userID string) ([]map[string]any, error) {
	rows, err := r.Db.Query(`
		SELECT DISTINCT u.id, u.firstname, u.lastname, u.profile_image
		FROM users u
		JOIN followers f ON (
			(f.follower_id = ? AND f.following_id = u.id AND f.status = 'accepted')
			OR
			(f.following_id = ? AND f.follower_id = u.id AND f.status = 'accepted')
		)
		WHERE u.id != ?
		ORDER BY u.firstname ASC, u.lastname ASC
	`, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []map[string]any{}
	for rows.Next() {
		var id, firstname, lastname, profileImage string
		if err := rows.Scan(&id, &firstname, &lastname, &profileImage); err != nil {
			return nil, err
		}
		users = append(users, map[string]any{
			"id":            id,
			"firstname":     firstname,
			"lastname":      lastname,
			"profile_image": profileImage,
		})
	}
	return users, rows.Err()
}
