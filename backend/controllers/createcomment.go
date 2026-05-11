package controllers

import (
	"net/http"

	"rtf/help"
	"rtf/models"
)

// create comments handler
func (c *Controller) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		help.RespondNotOK(w, "notallowed")
		return
	}

	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		help.RespondNotOK(w, "notallowed")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		help.RespondNotOK(w, "badrequest")
		return
	}

	var comment models.Comment
	comment.Content = r.FormValue("content")
	comment.UserID = userID

	postID := r.FormValue("post_id")
	comment.PostID = postID

	file, handler, err := r.FormFile("image_url")
	if err == nil {
		defer file.Close()

		if !help.IsPictureFormatCorrect(file, handler) {
			help.RespondNotOK(w, "badrequest")
			return
		}

		filename := help.SaveFile(file, handler.Filename)
		if filename != "" {
			comment.ImageURL = "/pics/" + filename
		}
	}

	if !help.IsvalidComment(comment) {
		help.RespondNotOK(w, "badrequest")
		return
	}

	if err := c.DB.PostExists(comment.PostID); err != nil {
		help.RespondNotOK(w, "badrequest")
		return
	}

	comment, err = c.DB.InsertCommentDB(comment)
	if err != nil {
		help.RespondNotOK(w, "badrequest")
		return
	}

	help.RespondOK(w, comment, "comment")

}
