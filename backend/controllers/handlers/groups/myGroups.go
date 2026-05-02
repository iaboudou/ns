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

func GetJoinedGroups(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.Context().Value("userID").(string)
	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")

	var rows *sql.Rows
	var err error

	if last == "" {
		rows, err = db.Query(`
			SELECT g.id, g.title, g.description, g.created_at
			FROM group_members gm
			JOIN groups g ON gm.group_id = g.id
			WHERE gm.user_id = ? AND gm.status = "accepted"
			ORDER BY g.created_at DESC, g.id DESC
			LIMIT 10
		`, userID)
	} else {
		normalized := strings.Replace(last, "T", " ", 1)
		normalized = strings.TrimSuffix(normalized, "Z")

		rows, err = db.Query(`
			SELECT g.id, g.title, g.description, g.created_at
			FROM group_members gm
			JOIN groups g ON gm.group_id = g.id
			WHERE gm.user_id = ? AND gm.status = "accepted"
			  AND (g.created_at < ? OR (g.created_at = ? AND g.id < ?))
			ORDER BY g.created_at DESC, g.id DESC
			LIMIT 10
		`, userID, normalized, normalized, lastID)
	}

	if err != nil {
		fmt.Println("error creating rows for geting groups: ", err)
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
			fmt.Println("error while getting users's group")
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
