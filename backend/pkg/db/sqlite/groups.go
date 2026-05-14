package sqlite

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"rtf/models"

	"github.com/gofrs/uuid/v5"
)

func InsertGroupMessage(db *sql.DB, groupID, userID, content string, createdAt int64) (string, error) {
	var isMember bool
	err := db.QueryRow(`SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?`, userID, groupID).Scan(&isMember)
	if err != nil {
		return "", err
	}

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

func InsertEventInDB(db *sql.DB, userID, groupID string, eventId uuid.UUID, event *models.Event) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
	INSERT INTO events (id, group_id, creator_id, title, description, event_date)
	VALUES (?, ?, ?, ?, ?, ?)
	`, eventId.String(), groupID, userID, event.Title, event.Description, event.Date)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO event_responses (event_id, user_id, status)
	VALUES (?, ?, ?)
	ON CONFLICT(event_id, user_id)
	DO UPDATE SET status = excluded.status;
	`, eventId.String(), userID, event.Vote)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func InsertGroupInDB(db *sql.DB, w http.ResponseWriter, userID string, groupId uuid.UUID, group *models.Group) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO groups (
			id,
			creator_id,
			title,
			description,
			image
		)
		VALUES (?, ?, ?, ?, ?)`,
		groupId.String(),
		userID,
		group.Title,
		group.Description,
		group.Image,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO group_members (
			group_id,
			user_id,
			role,
			status
		)
		VALUES (?, ?, 'creator', 'accepted')`,
		groupId.String(),
		userID,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func SelectUnknownGroupsInDB(db *sql.DB, userID, filter, last, lastID string) ([]map[string]any, error) {
	var rows *sql.Rows
	var err error

	pattern := "%" + filter + "%"

	if last == "" {
		if filter != "" {
			rows, err = db.Query(`
				SELECT g.id, g.title, g.description, g.created_at, g.image
				FROM groups g
				WHERE NOT EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = g.id AND gm.user_id = ?
				)
				AND g.title LIKE ?
				ORDER BY g.created_at DESC, g.id DESC
				LIMIT 8
			`, userID, pattern)
		} else {
			rows, err = db.Query(`
				SELECT g.id, g.title, g.description, g.created_at, g.image
				FROM groups g
				WHERE NOT EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = g.id AND gm.user_id = ?
				)
				ORDER BY g.created_at DESC, g.id DESC
				LIMIT 8
			`, userID)
		}
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		if filter != "" {
			rows, err = db.Query(`
				SELECT g.id, g.title, g.description, g.created_at, g.image
				FROM groups g
				WHERE NOT EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = g.id AND gm.user_id = ?
				)
				AND g.title LIKE ?
				AND (g.created_at < ? OR (g.created_at = ? AND g.id < ?))
				ORDER BY g.created_at DESC, g.id DESC
				LIMIT 8
			`, userID, pattern, normalized, normalized, lastID)
		} else {
			rows, err = db.Query(`
				SELECT g.id, g.title, g.description, g.created_at, g.image
				FROM groups g
				WHERE NOT EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = g.id AND gm.user_id = ?
				)
				AND (g.created_at < ? OR (g.created_at = ? AND g.id < ?))
				ORDER BY g.created_at DESC, g.id DESC
				LIMIT 8
			`, userID, normalized, normalized, lastID)
		}
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []map[string]any{}

	for rows.Next() {
		var id, title, description, image string
		var created_at time.Time

		err := rows.Scan(&id, &title, &description, &created_at, &image)
		if err != nil {
			return nil, err
		}

		groups = append(groups, map[string]any{
			"id":          id,
			"title":       title,
			"description": description,
			"created_at":  created_at,
			"img":         image,
		})
	}

	return groups, nil
}

func SelectEventsInDB(db *sql.DB, userID, groupID, last, lastID string) ([]map[string]any, error) {
	var rows *sql.Rows
	var err error

	if last == "" {
		rows, err = db.Query(`
			SELECT u.nickname, u.firstname, u.lastname, e.id, e.title, e.description, e.event_date, e.created_at
			FROM events e
			JOIN users u ON u.id = e.creator_id
			WHERE e.group_id = ?
			ORDER BY e.created_at DESC, e.id DESC
			LIMIT 10
		`, groupID)
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		rows, err = db.Query(`
			SELECT u.nickname, u.firstname, u.lastname, e.id, e.title, e.description, e.event_date, e.created_at
			FROM events e
			JOIN users u ON u.id = e.creator_id
			WHERE e.group_id = ?
			  AND (e.created_at < ? OR (e.created_at = ? AND e.id < ?))
			ORDER BY e.created_at DESC, e.id DESC
			LIMIT 10
		`, groupID, normalized, normalized, lastID)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []map[string]any{}

	for rows.Next() {
		var id, title, description, author_nickname, author_firstname, author_lastname, date string
		var created_at time.Time

		err := rows.Scan(&author_nickname, &author_firstname, &author_lastname, &id, &title, &description, &date, &created_at)
		if err != nil {
			return nil, err
		}

		var going, notgoing int
		var voted sql.NullString

		err = db.QueryRow(`
			SELECT
				COUNT(CASE WHEN status = 'going' THEN 1 END),
				COUNT(CASE WHEN status = 'notgoing' THEN 1 END),
				MAX(CASE WHEN user_id = ? THEN status END)
			FROM event_responses
			WHERE event_id = ?
		`, userID, id).Scan(&going, &notgoing, &voted)
		if err != nil {
			return nil, err
		}

		events = append(events, map[string]any{
			"id":             id,
			"title":          title,
			"description":    description,
			"date":           date,
			"created_at":     created_at,
			"goingCount":     going,
			"notGoingCount":  notgoing,
			"voted":          voted.String,
			"authorNickname": author_nickname,
			"authorName":     author_firstname + " " + author_lastname,
		})
	}

	return events, nil
}

func SelectGroupDescription(db *sql.DB, info *models.GroupeInfo, groupID, user_id string) error {
	var role *string

	err := db.QueryRow(`
	SELECT 
		g.title,
		g.description,
		g.image,
		(g.creator_id = ?) as is_creator,
		gm.role
	FROM groups g
	LEFT JOIN group_members gm 
		ON gm.group_id = g.id 
		AND gm.user_id = ?
		AND gm.status = 'accepted'
	WHERE g.id = ?
	`, user_id, user_id, groupID).Scan(
		&info.Title,
		&info.Description,
		&info.Image,
		&info.IsCreator,
		&role,
	)
	if err != nil {
		return err
	}

	if role == nil {
		return errors.New("not a member")
	}
	err = db.QueryRow(`SELECT COUNT(*) FROM group_members WHERE group_id = ?`, groupID).Scan(&info.MemberAmount)

	return err
}

func SelectGroupInvitesInDB(db *sql.DB, userID, last, lastID string) ([]map[string]any, error) {
	var rows *sql.Rows
	var err error
	if last == "" {
		rows, err = db.Query(`
        SELECT g.id, g.title, g.description, gm.joined_at, g.image, gm.invited_by
        FROM groups g
        JOIN group_members gm ON gm.group_id = g.id
        WHERE gm.user_id = ? AND gm.status = 'pending' AND gm.type = 'invite'
        ORDER BY gm.joined_at DESC, g.id DESC
        LIMIT 8
    `, userID)
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		rows, err = db.Query(`
        SELECT g.id, g.title, g.description, gm.joined_at, g.image, gm.invited_by
        FROM groups g
        JOIN group_members gm ON gm.group_id = g.id
        WHERE gm.user_id = ? AND gm.status = 'pending' AND gm.type = 'invite'
          AND (gm.joined_at < ? OR (gm.joined_at = ? AND g.id < ?))
        ORDER BY gm.joined_at DESC, g.id DESC
        LIMIT 8
    `, userID, normalized, normalized, lastID)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []map[string]any{}

	for rows.Next() {
		var id, title, description, image string
		var joined_at time.Time
		var invitedBy sql.NullString

		err := rows.Scan(&id, &title, &description, &joined_at, &image, &invitedBy)
		if err != nil {
			return nil, err
		}

		groups = append(groups, map[string]any{
			"id":          id,
			"title":       title,
			"description": description,
			"created_at":  joined_at,
			"img":         image,
			"invited_by":  invitedBy.String,
		})
	}

	return groups, err
}

func SelectGroupRequests(db *sql.DB, groupID, last, lastID string) ([]map[string]any, error) {
	var rows *sql.Rows
	var err error

	if last == "" {
		rows, err = db.Query(`
			SELECT u.id, u.nickname, u.firstname, u.lastname, gm.joined_at
			FROM group_members gm
			JOIN users u ON u.id = gm.user_id
			WHERE gm.group_id = ? AND gm.status = 'pending' AND gm.type = 'request'
			ORDER BY gm.joined_at DESC, u.id DESC
			LIMIT 10
		`, groupID)
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		rows, err = db.Query(`
			SELECT u.id, u.nickname, u.firstname, u.lastname, gm.joined_at
			FROM group_members gm
			JOIN users u ON u.id = gm.user_id
			WHERE gm.group_id = ? AND gm.status = 'pending' AND gm.type = 'request'
			  AND (gm.joined_at < ? OR (gm.joined_at = ? AND u.id < ?))
			ORDER BY gm.joined_at DESC, u.id DESC
			LIMIT 10
		`, groupID, normalized, normalized, lastID)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []map[string]any{}

	for rows.Next() {
		var id, nickname, firstname, lastname string
		var joined_at time.Time

		err := rows.Scan(&id, &nickname, &firstname, &lastname, &joined_at)
		if err != nil {
			return nil, err
		}

		users = append(users, map[string]any{
			"id":         id,
			"nickname":   nickname,
			"firstname":  firstname,
			"lastname":   lastname,
			"created_at": joined_at,
		})
	}

	return users, nil
}

func SelectFollowers(db *sql.DB, userId, groupID, last, lastID, search string) ([]map[string]any, error) {
	var rows *sql.Rows
	var err error

	pattern := "%" + search + "%"
	if last == "" {
		if search != "" {
			rows, err = db.Query(`
			SELECT u.id, u.nickname, u.firstname, u.lastname, u.profile_image, u.created_at
			FROM users u
			JOIN followers f ON f.following_id = u.id
			WHERE f.follower_id = ?
			  AND f.status = 'accepted'
			  AND u.id != ?
			  AND u.id NOT IN (
			      SELECT user_id FROM group_members WHERE group_id = ?
			  )
			  AND (
			    u.nickname LIKE ? OR
			    u.firstname LIKE ? OR
			    u.lastname LIKE ? OR
			    u.firstname || ' ' || u.lastname LIKE ?
			  )
			ORDER BY u.created_at DESC, u.id DESC
			LIMIT 10
		`, userId, userId, groupID, pattern, pattern, pattern, pattern)
		} else {
			rows, err = db.Query(`
			SELECT u.id, u.nickname, u.firstname, u.lastname, u.profile_image, u.created_at
			FROM users u
			JOIN followers f ON f.following_id = u.id
			WHERE f.follower_id = ?
			  AND f.status = 'accepted'
			  AND u.id != ?
			  AND u.id NOT IN (
			      SELECT user_id FROM group_members WHERE group_id = ?
			  )
			ORDER BY u.created_at DESC, u.id DESC
			LIMIT 10
		`, userId, userId, groupID)
		}
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		rows, err = db.Query(`
		SELECT u.id, u.nickname, u.firstname, u.lastname, u.profile_image, u.created_at
		FROM users u
		JOIN followers f ON f.following_id = u.id
		WHERE f.follower_id = ?
		  AND f.status = 'accepted'
		  AND (u.created_at < ? OR (u.created_at = ? AND u.id < ?))
		  AND u.id != ?
		  AND u.id NOT IN (
		      SELECT user_id FROM group_members WHERE group_id = ?
		  )
		ORDER BY u.created_at DESC, u.id DESC
		LIMIT 10
	`, userId, normalized, normalized, lastID, userId, groupID)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []map[string]any{}

	for rows.Next() {
		var id, nickname, firstname, lastname, profile_image string
		var created_at time.Time

		err := rows.Scan(&id, &nickname, &firstname, &lastname, &profile_image, &created_at)
		if err != nil {
			return nil, err
		}

		users = append(users, map[string]any{
			"id":            id,
			"nickname":      nickname,
			"firstname":     firstname,
			"lastname":      lastname,
			"fullname":      firstname + " " + lastname,
			"profile_image": profile_image,
			"created_at":    created_at,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func InsertGroupMember(db *sql.DB, userID, groupID, currentUserID, Type string) error {
	var err error
	if Type == "invite" {
		_, err = db.Exec(
			`INSERT INTO group_members (group_id, user_id, type, invited_by) VALUES (?, ?, ?, ?)`,
			groupID, userID, Type, currentUserID,
		)
	} else {
		_, err = db.Exec(
			`INSERT INTO group_members (group_id, user_id, type) VALUES (?, ?, ?)`,
			groupID, userID, Type,
		)
	}

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("already member")
		}

		return err
	}

	return nil
}

func UpdateRequest(db *sql.DB, request *models.Request, groupID string) error {
	var res sql.Result
	var err error

	switch request.Decicion {
	case "rejected":
		res, err = db.Exec(`DELETE FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, request.Sender)
	case "accepted":
		res, err = db.Exec(`UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ?`, request.Decicion, groupID, request.Sender)
	default:
		return errors.New("unknown decision")
	}

	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("not found")
	}

	return nil
}

func RemoveNotif(db *sql.DB, request *models.Request, userID, groupID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if request.InvitedBy != "" {
		_, err = tx.Exec(`
			DELETE FROM notification_users
			WHERE user_id = ?
			AND notification_id IN (
				SELECT id FROM notifications
				WHERE group_id = ?
				AND type = 'group_invite'
				AND sender_id = ?
			)
		`, userID, groupID, request.InvitedBy)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM notifications
			WHERE group_id = ?
			AND type = 'group_invite'
			AND sender_id = ?
		`, groupID, request.InvitedBy)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec(`
			DELETE FROM notification_users
			WHERE user_id = ?
			AND notification_id IN (
				SELECT id FROM notifications
				WHERE group_id = ?
				AND type = 'group_request'
				AND sender_id = ?
			)
		`, userID, groupID, request.Sender)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM notifications
			WHERE group_id = ?
			AND type = 'group_request'
			AND sender_id = ?
		`, groupID, request.Sender)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func SelectMyGroups(db *sql.DB, userID, last, lastID string) ([]map[string]any, error) {
	var rows *sql.Rows
	var err error

	if last == "" {
		rows, err = db.Query(`
			SELECT g.id, g.title, g.description, g.created_at, g.image
			FROM group_members gm
			JOIN groups g ON gm.group_id = g.id
			WHERE gm.user_id = ? AND gm.status = "accepted"
			ORDER BY g.created_at DESC, g.id DESC
			LIMIT 8
		`, userID)
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		rows, err = db.Query(`
			SELECT g.id, g.title, g.description, g.created_at, g.image
			FROM group_members gm
			JOIN groups g ON gm.group_id = g.id
			WHERE gm.user_id = ? AND gm.status = "accepted"
			  AND (g.created_at < ? OR (g.created_at = ? AND g.id < ?))
			ORDER BY g.created_at DESC, g.id DESC
			LIMIT 8
		`, userID, normalized, normalized, lastID)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []map[string]any{}

	for rows.Next() {
		var id, title, description, image string
		var created_at time.Time

		err := rows.Scan(&id, &title, &description, &created_at, &image)
		if err != nil {
			return nil, err
		}

		groups = append(groups, map[string]any{
			"id":          id,
			"title":       title,
			"description": description,
			"created_at":  created_at,
			"img":         image,
		})
	}

	return groups, nil
}

func InsertVoteInDB(db *sql.DB, userID, eventID, vote string) error {
	_, err := db.Exec(`
			INSERT INTO event_responses (event_id, user_id, status)
			VALUES (?, ?, ?)
			ON CONFLICT(event_id, user_id)
			DO UPDATE SET status = excluded.status;
		`, eventID, userID, vote)

	return err
}
