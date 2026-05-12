package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"rtf/help"
	"rtf/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID, userId string) {
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

	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")
	search := r.URL.Query().Get("search")

	var rows *sql.Rows

	pattern := "%" + search + "%"
	if last == "" {
		if search != "" {
			rows, err = db.Query(`
			SELECT u.id, u.nickname, u.firstname, u.lastname, u.profile_image, u.created_at
			FROM users u
			JOIN followers f ON f.follower_id = u.id
			WHERE f.following_id = ?
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
			JOIN followers f ON f.follower_id = u.id
			WHERE f.following_id = ?
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
		JOIN followers f ON f.follower_id = u.id
		WHERE f.following_id = ?
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
		help.RespondServerError(w)
		return
	}

	defer rows.Close()

	users := []map[string]any{}

	for rows.Next() {
		var id, nickname, firstname, lastname, profile_image string
		var created_at time.Time

		err := rows.Scan(&id, &nickname, &firstname, &lastname, &profile_image, &created_at)
		if err != nil {
			help.RespondServerError(w)
			return
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
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: 200,
		Data: users,
	})
}
