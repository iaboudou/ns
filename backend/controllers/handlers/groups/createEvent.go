package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"rtf/help"
	"rtf/models"

	"rtf/pkg/db/sqlite"

	"github.com/gofrs/uuid/v5"
)

func CreateEvent(w http.ResponseWriter, r *http.Request, db *sqlite.Repo, hub *models.Hub, groupID, userID string) {
	var event models.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid body",
		})
		return
	}

	mistake := ValidateEvent(&event)
	if mistake != "" {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: mistake,
		})
		return
	}

	eventId, err := uuid.NewV4()
	if err != nil {
		help.RespondServerError(w)
		return
	}

	tx, err := db.Db.Begin()
	if err != nil {
		help.RespondServerError(w)
		return
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
	INSERT INTO events (id, group_id, creator_id, title, description, event_date)
	VALUES (?, ?, ?, ?, ?, ?)
	`, eventId.String(), groupID, userID, event.Title, event.Description, event.Date)
	if err != nil {
		tx.Rollback()
		help.RespondServerError(w)
		return
	}

	_, err = tx.Exec(`
	INSERT INTO event_responses (event_id, user_id, status)
	VALUES (?, ?, ?)
	ON CONFLICT(event_id, user_id)
	DO UPDATE SET status = excluded.status;
	`, eventId.String(), userID, event.Vote)
	if err != nil {
		tx.Rollback()
		help.RespondServerError(w)
		return
	}

	err = tx.Commit()
	if err != nil {
		help.RespondServerError(w)
		return
	}

	hub.Notif <- models.Notification{
		SenderID: userID,
		Type:     "event",
		GroupID:  groupID,
	}

	var author_nickname, author_firstname, author_lastname string
	err = db.Db.QueryRow(`SELECT nickname, firstname, lastname FROM users WHERE
	id = ?`, userID).Scan(&author_nickname, &author_firstname, &author_lastname)
	if err != nil {
		help.RespondServerError(w)
		return
	}
	var goingCount, notGoingCount int
	if event.Vote == "going" {
		goingCount++
	} else {
		notGoingCount++
	}

	help.Respond(w, &models.Response{
		Code: http.StatusCreated,
		Data: map[string]any{
			"id":             eventId.String(),
			"title":          event.Title,
			"description":    event.Description,
			"date":           event.Date,
			"goingCount":     goingCount,
			"notGoingCount":  notGoingCount,
			"voted":          event.Vote,
			"authorNickname": author_nickname,
			"authorName":     author_firstname + " " + author_lastname,
		},
	})
}

func ValidateEvent(e *models.Event) string {
	e.Title = strings.TrimSpace(e.Title)
	e.Description = strings.TrimSpace(e.Description)
	e.Vote = strings.ToLower(strings.TrimSpace(e.Vote))

	if e.Title == "" || e.Description == "" {
		return "Please fill all the fields"
	}

	if len(e.Title) < 3 {
		return "Title is too short"
	}

	if len(e.Title) > 80 {
		return "title is too long"
	}

	if len(e.Description) < 10 {
		return "description is too short"
	}

	if len(e.Description) > 500 {
		return "description is too long"
	}

	if e.Vote != "going" && e.Vote != "notgoing" {
		return "you can only chose to go or not to go"
	}

	date, err := time.Parse(time.RFC3339, e.Date)
	if err != nil {
		return "invalid time format"
	}

	minTime := time.Now().Add(5 * time.Minute)

	if !date.After(minTime) {
		return "an event must be at least in 5 minutes"
	}

	return ""
}
