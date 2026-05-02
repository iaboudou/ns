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

func GetUnknownGroups(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.Context().Value("userID").(string)
	filter := r.URL.Query().Get("search")
	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")

	var rows *sql.Rows
	var err error

	pattern := "%" + filter + "%"

	if last == "" {
		if filter != "" {
			rows, err = db.Query(`
				SELECT g.id, g.title, g.description, g.created_at
				FROM groups g
				WHERE NOT EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = g.id AND gm.user_id = ?
				)
				AND g.title LIKE ?
				ORDER BY g.created_at DESC, g.id DESC
				LIMIT 10
			`, userID, pattern)
		} else {
			rows, err = db.Query(`
				SELECT g.id, g.title, g.description, g.created_at
				FROM groups g
				WHERE NOT EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = g.id AND gm.user_id = ?
				)
				ORDER BY g.created_at DESC, g.id DESC
				LIMIT 10
			`, userID)
		}
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		if filter != "" {
			rows, err = db.Query(`
				SELECT g.id, g.title, g.description, g.created_at
				FROM groups g
				WHERE NOT EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = g.id AND gm.user_id = ?
				)
				AND g.title LIKE ?
				AND (g.created_at < ? OR (g.created_at = ? AND g.id < ?))
				ORDER BY g.created_at DESC, g.id DESC
				LIMIT 10
			`, userID, pattern, normalized, normalized, lastID)
		} else {
			rows, err = db.Query(`
				SELECT g.id, g.title, g.description, g.created_at
				FROM groups g
				WHERE NOT EXISTS (
					SELECT 1 FROM group_members gm
					WHERE gm.group_id = g.id AND gm.user_id = ?
				)
				AND (g.created_at < ? OR (g.created_at = ? AND g.id < ?))
				ORDER BY g.created_at DESC, g.id DESC
				LIMIT 10
			`, userID, normalized, normalized, lastID)
		}
	}

	if err != nil {
		fmt.Println("error creating rows for group suggestion: ", err)
		help.RespondServerError(w)
		return
	}

	defer rows.Close()

	groups := []map[string]any{}

	for rows.Next() {
		var id, title, description string
		var created_at time.Time

		err := rows.Scan(&id, &title, &description, &created_at)
		if err != nil {
			fmt.Println("error while getting groups suggestions :", err)
			help.RespondServerError(w)
			return
		}

		groups = append(groups, map[string]any{
			"id":          id,
			"title":       title,
			"description": description,
			"created_at":  created_at,
		})
	}

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: groups,
	})
}
