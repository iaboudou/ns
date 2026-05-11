package controllers

import (
	"net/http"
	"strings"

	"rtf/models"
	"rtf/help"
)

// handle create post
func (c *Controller) CreatePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// get the user ID

	if r.Method != http.MethodPost {
		help.RespondNotOK(w, "notallowed")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		help.RespondNotOK(w, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		help.RespondNotOK(w, "badrequest")
		return
	}

	var post models.Post
	post.UserID = userID
	post.Content = strings.TrimSpace(r.FormValue("text"))
	post.Privacy = strings.TrimSpace(r.FormValue("privacy"))

	if post.Privacy != "public" && post.Privacy != "private" && post.Privacy != "followers" {
		help.RespondNotOK(w, "badrequest")
		return
	}

	post.GroupID = strings.TrimSpace(r.FormValue("group_id"))

	if post.Privacy == "private" {
		post.Alloweduserscreate = r.FormValue("allowed_users")
	} else {
		post.AllowedUsers = nil
	}

	// handle the file upload
	f, h, er := r.FormFile("Image")
	if er == nil {
		defer f.Close()

		if !help.IsPictureFormatCorrect(f, h) {
			help.RespondNotOK(w, "badrequest")
			return
		}

		defer f.Close()
		filename := help.SaveFile(f, h.Filename)
		if filename != "" {
			post.ImageURL = "/pics/" + filename
		}
	}

	// check if the post content is correct
	err := help.ArePostInfosCorrect(post)
	if err != nil {
		help.RespondNotOK(w, "badrequest")
		return
	}
	// insert the post into the DB
	post, er = c.DB.InsertPostDB(userID, post)
	if er != nil {
		help.RespondNotOK(w, "badrequest")
		return
	}
	post.NumberOfComments = 0

	help.RespondOK(w, post, "post")
}
