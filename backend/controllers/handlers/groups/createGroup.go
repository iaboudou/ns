package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"rtf/help"
	"rtf/models"

	"github.com/gofrs/uuid/v5"
)

func CreateGroup(w http.ResponseWriter, r *http.Request, db *sql.DB, userID string) {
	var group models.Group

	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid body",
		})
		return
	}

	mistake := ValidateGroup(&group)
	if mistake != "" {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: mistake,
		})
		return
	}

	groupId, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error generating group id: ", err)
		help.RespondServerError(w)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("error generating transaction for group creation: ", err)
		help.RespondServerError(w)
		return
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
			INSERT INTO groups (id, creator_id, title, description )
			Values (?, ?, ?, ?)`, groupId.String(), userID, group.Title, group.Description)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			help.Respond(w, &models.Response{
				Code:    http.StatusConflict,
				Message: "This title is already used, please select another one",
			})
			return
		}

		fmt.Println("error inserting group in group table: ", err)
		help.RespondServerError(w)
		return
	}

	_, err = tx.Exec(`
			INSERT INTO group_members (group_id, user_id, role, status)
			VALUES (?, ?, 'creator', 'accepted')`, groupId.String(), userID)
	if err != nil {
		fmt.Println("error inserting creator group in group_members table: ", err)
		help.RespondServerError(w)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("error commiting transactions: ", err)
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusCreated,
		Data: map[string]string{
			"id":          groupId.String(),
			"title":       group.Title,
			"description": group.Description,
		},
	})
}

func ValidateGroup(g *models.Group) string {
	g.Title = strings.TrimSpace(g.Title)
	g.Description = strings.TrimSpace(g.Description)

	if g.Title == "" || g.Description == "" {
		return "Please fill all the fields"
	}

	if len(g.Title) < 3 {
		return "Title is too short"
	}

	if len(g.Title) > 60 {
		return "title is too long"
	}

	if len(g.Description) < 20 {
		return "description is too short"
	}

	if len(g.Description) > 300 {
		return "description is too long"
	}
	return ""
}
