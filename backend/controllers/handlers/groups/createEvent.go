package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"rtf/help"
	"rtf/models"

	"rtf/pkg/db/sqlite"

	"github.com/gofrs/uuid/v5"
)

func CreateEvent(w http.ResponseWriter, r *http.Request, db *sql.DB, hub *models.Hub, groupID, userID string) {
	var exists int
	err := db.QueryRow(`SELECT 1 FROM groups WHERE id = ?`, groupID).Scan(&exists) // check group existence
	if err == sql.ErrNoRows {
		help.RespondNotFound(w, "This group doesn't exist")
		return
	}

	if err != nil {
		help.RespondServerError(w)
		return
	}

	var isMember bool
	err = db.QueryRow(`SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?`, userID, groupID).Scan(&isMember)
	if err != nil {
		if err == sql.ErrNoRows {
			help.Respond(w, &models.Response{
				Code:    http.StatusBadRequest,
				Message: "you are not member of the group",
			})
			return
		}

		help.RespondServerError(w)
		return
	}

	var event models.Event

	err = json.NewDecoder(r.Body).Decode(&event)
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

	hub.Notif <- models.Notification{
		SenderID: userID,
		Type:     "event",
		GroupID:  groupID,
	}

	err = sqlite.InsertEventInDB(db, userID, groupID, eventId, &event)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	author, err := sqlite.GetUserByID(db, userID)
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
			"authorNickname": author.Nickname,
			"authorName":     author.Firstname + " " + author.Lastname,
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
