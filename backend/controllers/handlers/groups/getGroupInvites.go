package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"rtf/help"
	"rtf/models"
)

func GetGroupInvites(w http.ResponseWriter, r *http.Request, db *sql.DB, userID string) {
	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")

	var rows *sql.Rows
	var err error

	if last == "" {
		rows, err = db.Query(`
			SELECT g.id, g.title, g.description, gm.joined_at, g.image
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
			SELECT g.id, g.title, g.description, gm.joined_at, g.image
			FROM groups g
			JOIN group_members gm ON gm.group_id = g.id
			WHERE gm.user_id = ? AND gm.status = 'pending' AND gm.type = 'invite'
			  AND (gm.joined_at < ? OR (gm.joined_at = ? AND g.id < ?))
			ORDER BY gm.joined_at DESC, g.id DESC
			LIMIT 8
		`, userID, normalized, normalized, lastID)
	}

	if err != nil {
		help.RespondServerError(w)
		return
	}

	defer rows.Close()

	groups := []map[string]any{}

	for rows.Next() {
		var id, title, description, image string
		var joined_at time.Time

		err := rows.Scan(&id, &title, &description, &joined_at, &image)
		if err != nil {
			help.RespondServerError(w)
			return
		}

		groups = append(groups, map[string]any{
			"id":          id,
			"title":       title,
			"description": description,
			"created_at":  joined_at,
			"img":         image,
		})
	}

	help.Respond(w, &models.Response{
		Code: 200,
		Data: groups,
	})
}
