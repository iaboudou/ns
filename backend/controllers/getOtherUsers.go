package controllers

import (
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

	query := `
		SELECT 
			u.id, u.nickname, u.firstname, u.lastname, u.profile_image, u.created_at,
			(SELECT COUNT(*) FROM messages 
			 WHERE sender_id = u.id AND receiver_id = ? AND is_not_read = 1) as unread_count
		FROM users u
		WHERE u.id IN (
			SELECT following_id FROM followers WHERE follower_id = ? AND status = 'accepted'
			UNION
			SELECT follower_id FROM followers WHERE following_id = ? AND status = 'accepted'
		)`

	args := []any{userID, userID, userID}

	if search != "" {
		pattern := "%" + search + "%"
		query += " AND (u.nickname LIKE ? OR u.firstname LIKE ? OR u.lastname LIKE ? OR (u.firstname || ' ' || u.lastname) LIKE ? OR u.id = ?)"
		args = append(args, pattern, pattern, pattern, pattern, search)
	}

	if last != "" {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")
		query += " AND (u.created_at < ? OR (u.created_at = ? AND u.id < ?))"
		args = append(args, normalized, normalized, lastID)
	}

	query += " ORDER BY u.created_at DESC, u.id DESC LIMIT 10"

	rows, err := db.Query(query, args...)
	if err != nil {
		help.RespondServerError(w)
		return
	}
	defer rows.Close()

	users := []map[string]any{}
	for rows.Next() {
		var u struct {
			ID           string
			Nickname     string
			Firstname    string
			Lastname     string
			ProfileImage string
			CreatedAt    time.Time
			UnreadCount  int
		}
		if err := rows.Scan(&u.ID, &u.Nickname, &u.Firstname, &u.Lastname, &u.ProfileImage, &u.CreatedAt, &u.UnreadCount); err != nil {
			continue
		}

		users = append(users, map[string]any{
			"id":            u.ID,
			"nickname":      u.Nickname,
			"firstname":     u.Firstname,
			"lastname":      u.Lastname,
			"profile_image": u.ProfileImage,
			"created_at":    u.CreatedAt,
			"unread_count":  u.UnreadCount,
		})
	}

	help.Respond(w, &models.Response{
		Code: 200,
		Data: users,
	})
}
