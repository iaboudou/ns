package handlers

import (
	"database/sql"
	"fmt"
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
				SELECT id, nickname, firstname, lastname, profile_image, created_at
				FROM users
				WHERE (nickname LIKE ? OR firstname LIKE ? OR lastname LIKE ?)
				  AND id != ?
				  AND id NOT IN (SELECT user_id FROM group_members WHERE group_id = ?)
				ORDER BY created_at DESC, id DESC
				LIMIT 10
			`, pattern, pattern, pattern, userId, groupID)
		} else {
			rows, err = db.Query(`
				SELECT id, nickname, firstname, lastname, profile_image, created_at
				FROM users
				WHERE id != ?
				  AND id NOT IN (SELECT user_id FROM group_members WHERE group_id = ?)
				ORDER BY created_at DESC, id DESC
				LIMIT 10
			`, userId, groupID)
		}
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		rows, err = db.Query(`
			SELECT id, nickname, firstname, lastname, profile_image, created_at
			FROM users
			WHERE (created_at < ? OR (created_at = ? AND id < ?))
			  AND id != ?
			  AND id NOT IN (SELECT user_id FROM group_members WHERE group_id = ?)
			ORDER BY created_at DESC, id DESC
			LIMIT 10
		`, normalized, normalized, lastID, userId, groupID)
	}

	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
			help.RespondServerError(w)
			return
		}

		users = append(users, map[string]any{
			"id":            id,
			"nickname":      nickname,
			"firstname":     firstname,
			"lastname":      lastname,
			"profile_image": profile_image,
			"created_at":    created_at,
		})
	}

	help.Respond(w, &models.Response{
		Code: 200,
		Data: users,
	})
}
