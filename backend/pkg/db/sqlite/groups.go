package sqlite

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid/v5"
)

// 
func InsertGroupMessage(db *sql.DB, groupID, userID, content string, createdAt int64) (string, error) {
	messageID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(`
		INSERT INTO group_messages (id, group_id, user_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, messageID.String(), groupID, userID, content, createdAt)

	return messageID.String(), err
}

//
func SelectOldGroupMessages(db *sql.DB, groupID string, lastTime int64) ([]map[string]any, error) {
	if lastTime == 0 {
		lastTime = time.Now().UnixMilli()
	}

	rows, err := db.Query(`
		SELECT gm.id, gm.group_id, gm.user_id, gm.content, gm.created_at,
		       u.nickname, u.firstname, u.lastname, u.profile_image
		FROM group_messages gm
		JOIN users u ON u.id = gm.user_id
		WHERE gm.group_id = ? AND gm.created_at < ?
		ORDER BY gm.created_at DESC
		LIMIT 10
	`, groupID, lastTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []map[string]any{}
	for rows.Next() {
		var id, gId, uId, content string
		var createdAt int64
		var nickname, firstname, lastname, profile string

		if err := rows.Scan(&id, &gId, &uId, &content, &createdAt, &nickname, &firstname, &lastname, &profile); err != nil {
			continue
		}

		messages = append(messages, map[string]any{
			"id":              id,
			"group_id":        gId,
			"sender_id":       uId,
			"content":         content,
			"created_at":      createdAt,
			"sender_nickname": nickname,
			"sender_fullname": firstname + " " + lastname,
			"sender_profile":  profile,
			"is_group":        true,
		})
	}

	return messages, nil
}

//
func GetGroupMembers(db *sql.DB, groupID string) ([]string, error) {
	rows, err := db.Query(`SELECT user_id FROM group_members WHERE group_id = ? AND status = 'accepted'`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			continue
		}
		members = append(members, userID)
	}
	return members, nil
}
