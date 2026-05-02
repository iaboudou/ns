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

func GetGroupRequests(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID string) {
	userID := r.Context().Value("userID").(string)

	var creatorID string
	err := db.QueryRow(`SELECT creator_id FROM groups WHERE id = ?`, groupID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		help.RespondNotFound(w, "This group doesn't exist")
		return
	}
	if err != nil {
		help.RespondServerError(w)
		return
	}

	if userID != creatorID {
		help.Respond(w, &models.Response{
			Code:    http.StatusForbidden,
			Message: "not creator",
		})
		return
	}

	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")

	var rows *sql.Rows

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
		fmt.Println("error generating query to retrieve group request", err)
		help.RespondServerError(w)
		return
	}

	defer rows.Close()

	users := []map[string]any{}

	for rows.Next() {
		var id, nickname, firstname, lastname string
		var joined_at time.Time

		err := rows.Scan(&id, &nickname, &firstname, &lastname, &joined_at)
		if err != nil {
			fmt.Println("error in scan while retrieving group request", err)
			help.RespondServerError(w)
			return
		}

		users = append(users, map[string]any{
			"id":         id,
			"nickname":   nickname,
			"firstname":  firstname,
			"lastname":   lastname,
			"created_at": joined_at,
		})
	}

	help.Respond(w, &models.Response{Code: http.StatusOK, Data: users})
}
