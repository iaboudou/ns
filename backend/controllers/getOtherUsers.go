package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"rtf/help"
	"rtf/models"
)

func (c *Controller) GetOtherUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		help.RespondMethodNotAllowed(w)
		return
	}

	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")
	search := r.URL.Query().Get("search")
	db := c.DB.Db
	userID := r.Context().Value("userID").(string)

	// 1. get ids of people i have an accepted relationship with
	rows, err := db.Query(`
		SELECT following_id FROM followers WHERE follower_id = ? AND status = 'accepted'
		UNION
		SELECT follower_id FROM followers WHERE following_id = ? AND status = 'accepted'
	`, userID, userID)
	if err != nil {
		help.RespondServerError(w)
		return
	}
	defer rows.Close()

	var friendIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			friendIDs = append(friendIDs, id)
		}
	}

	if len(friendIDs) == 0 {
		help.Respond(w, &models.Response{Code: 200, Data: []any{}})
		return
	}

	placeholders := make([]string, len(friendIDs))
	args := make([]any, len(friendIDs))
	for i, id := range friendIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, nickname, firstname, lastname, profile_image, created_at
		FROM users
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	pattern := "%" + search + "%"
	if search != "" {
		query += " AND (nickname LIKE ? OR firstname LIKE ? OR lastname LIKE ? OR id = ?)"
		args = append(args, pattern, pattern, pattern, search)
	}

	if last != "" {
		// 
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")
		query += " AND (created_at < ? OR (created_at = ? AND id < ?))"
		args = append(args, normalized, normalized, lastID)
	}

	query += " ORDER BY created_at DESC, id DESC LIMIT 10"

	userRows, err := db.Query(query, args...)
	if err != nil {
		help.RespondServerError(w)
		return
	}
	defer userRows.Close()

	users := []map[string]any{}

	// 
	for userRows.Next() {
		var u struct {
			ID           string
			Nickname     string
			Firstname    string
			Lastname     string
			ProfileImage string
			CreatedAt    time.Time
		}
		if err := userRows.Scan(&u.ID, &u.Nickname, &u.Firstname, &u.Lastname, &u.ProfileImage, &u.CreatedAt); err != nil {
			continue
		}
		
		var lastMsg int64
		var unreadCount int

		// Get last message timestamp
		db.QueryRow(`
			SELECT IFNULL(MAX(created_at), 0) FROM messages
			WHERE (sender_id = ? AND receiver_id = ?)
			   OR (sender_id = ? AND receiver_id = ?)
		`, userID, u.ID, u.ID, userID).Scan(&lastMsg)

		// Get unread count 
		db.QueryRow(`
			SELECT COUNT(*) FROM messages
			WHERE sender_id = ? AND receiver_id = ? AND is_not_read = 1
		`, u.ID, userID).Scan(&unreadCount)

		users = append(users, map[string]any{
			"id":            u.ID,
			"nickname":      u.Nickname,
			"firstname":     u.Firstname,
			"lastname":      u.Lastname,
			"profile_image": u.ProfileImage,
			"created_at":    u.CreatedAt,
			"last_msg":      lastMsg,
			"unread_count":  unreadCount,
		})
	}

	help.Respond(w, &models.Response{
		Code: 200,
		Data: users,
	})
}
