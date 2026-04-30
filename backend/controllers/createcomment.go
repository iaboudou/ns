package controllers

import (
	"net/http"

	"rtf/models"
	"rtf/pkg"
)

// create comments handler
func (c *Controller) CreateComment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		pkg.RespondNotOK(w, "notallowed")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		pkg.RespondNotOK(w, "badrequest")
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

		if !pkg.IsPictureFormatCorrect(file, handler) {
			pkg.RespondNotOK(w, "badrequest")
			return
		}

		filename := pkg.SaveFile(file, handler.Filename)
		if filename != "" {
			comment.ImageURL = "/pics/" + filename
		}
	}

	if !pkg.IsvalidComment(comment) {
		pkg.RespondNotOK(w, "badrequest")
		return
	}

	if err := c.DB.PostExists(comment.PostID); err != nil {
		pkg.RespondNotOK(w, "badrequest")
		return
	}

	comment, err = c.DB.InsertCommentDB(comment)
	if err != nil {
		pkg.RespondNotOK(w, "badrequest")
		return
	}

	pkg.RespondOK(w, comment, "comment")

}
