package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"rtf/help"
	"rtf/models"

	"github.com/gofrs/uuid/v5"
)

func CreateEvent(w http.ResponseWriter, r *http.Request, db *sql.DB, groupID, userID string) {
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
		fmt.Println("error generating event id: ", err)
		help.RespondServerError(w)
		return
	}

	_, err = db.Exec(`INSERT INTO events (id, group_id,creator_id, title, description, event_date)
	VALUES (?, ?, ?, ?, ?, ?)`, eventId.String(), groupID, userID, event.Title, event.Description, event.Date)
	if err != nil {
		fmt.Println("error inserting event in db :", err)
		help.RespondServerError(w)
		return
	}

	var author_nickname, author_firstname, author_lastname string
	err = db.QueryRow(`SELECT nickname, firstname, lastname FROM users WHERE
	id = ?`, userID).Scan(&author_nickname, &author_firstname, &author_lastname)
	if err != nil {
		fmt.Println("error getting user name creating event: ", err)
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusCreated,
		Data: map[string]any{
			"id":             eventId.String(),
			"title":          event.Title,
			"description":    event.Description,
			"date":           event.Date,
			"goingCount":     0,
			"notGoingCount":  0,
			"voted":          "",
			"authorNickname": author_nickname,
			"authorName":     author_firstname + " " + author_lastname,
		},
	})
}

func ValidateEvent(e *models.Event) string {
	e.Title = strings.TrimSpace(e.Title)
	e.Description = strings.TrimSpace(e.Description)

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
