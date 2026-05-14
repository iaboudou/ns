package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"

	"github.com/gofrs/uuid/v5"
)

func CreateGroup(w http.ResponseWriter, r *http.Request, db *sql.DB, userID string) {
	var group models.Group

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid form data",
		})
		return
	}

	group.Title = r.FormValue("title")
	group.Description = r.FormValue("description")

	file, header, err := r.FormFile("image")
	if err == nil && header.Size != 0 {
		defer file.Close()
		if !help.IsPictureFormatCorrect(file, header) {
			help.Respond(w, &models.Response{
				Code:    http.StatusBadRequest,
				Message: "invalid image",
			})
			return
		}

		filename := help.SaveFile(file, header.Filename)
		if filename == "" {
			help.RespondServerError(w)
			return
		}
		group.Image = filename
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
		help.RespondServerError(w)
		return
	}

	err = sqlite.InsertGroupInDB(db, w, userID, groupId, &group)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			help.Respond(w, &models.Response{
				Code:    http.StatusConflict,
				Message: "This title is already used, please select another one",
			})
			return
		}

		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusCreated,
		Data: map[string]string{
			"id":          groupId.String(),
			"title":       group.Title,
			"description": group.Description,
			"img":         group.Image,
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
