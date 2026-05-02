package controllers

import (
	"database/sql"
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

	var (
		rows *sql.Rows
		err  error
	)

	pattern := "%" + search + "%"
	// "last" emty mean get the users in the first time else do the pagination
	if last == "" {
		// if search is emty do normal pagination else do the pagination with search
		if search != "" {
			rows, err = db.Query(`
				SELECT id, nickname, firstname, lastname, profile_image, created_at
				FROM users
				WHERE (nickname LIKE ? OR firstname LIKE ? OR lastname LIKE ?)
				AND id != ?
				ORDER BY created_at DESC, id DESC
				LIMIT 10
			`, pattern, pattern, pattern, userID)
		} else {
			rows, err = db.Query(`
				SELECT id, nickname, firstname, lastname, profile_image, created_at
				FROM users
				WHERE id != ?
				ORDER BY created_at DESC, id DESC
				LIMIT 10
			`, userID)
		}
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")
		if search != "" {
			rows, err = db.Query(`
            SELECT id, nickname, firstname, lastname, profile_image, created_at
            FROM users
            WHERE (created_at < ? OR (created_at = ? AND id < ?))
              AND id != ?
              AND (nickname LIKE ? OR firstname LIKE ? OR lastname LIKE ?)
            ORDER BY created_at DESC, id DESC
            LIMIT 10
        `, normalized, normalized, lastID, userID, pattern, pattern, pattern)
		} else {
			rows, err = db.Query(`
            SELECT id, nickname, firstname, lastname, profile_image, created_at
            FROM users
            WHERE (created_at < ? OR (created_at = ? AND id < ?))
              AND id != ?
            ORDER BY created_at DESC, id DESC
            LIMIT 10
        `, normalized, normalized, lastID, userID)
		}
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
