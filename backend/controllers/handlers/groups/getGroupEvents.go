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

func GetGroupEvents(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID, userID string) {
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

	var rows *sql.Rows

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
		fmt.Println("error creating rows to retrieve event: ", err)
		help.RespondServerError(w)
		return
	}

	defer rows.Close()

	events := []map[string]any{}

	for rows.Next() {
		var id, title, description, author_nickname, author_firstname, author_lastname, date string
		var created_at time.Time

		err := rows.Scan(&author_nickname, &author_firstname, &author_lastname, &id, &title, &description, &date, &created_at)
		if err != nil {
			fmt.Println("error while retrieving event basic data:", err)
			help.RespondServerError(w)
			return
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
			fmt.Println("error while retrieving event responses: ", err)
			help.RespondServerError(w)
			return
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

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: events,
	})
}
