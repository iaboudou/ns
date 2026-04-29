package controllers

import (
	"fmt"
	"net/http"
	"rtf/models"
	"rtf/pkg"
	"strings"
)

// handle create post
func (c *Controller) CreatePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// get the user ID

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		pkg.RespondNotOK(w, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		pkg.RespondNotOK(w, "badrequest")
		return
	}

	var post models.Post
	post.UserID = userID
	post.Content = strings.TrimSpace(r.FormValue("text"))
	post.Privacy = strings.TrimSpace(r.FormValue("privacy"))
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

		if !pkg.IsPictureFormatCorrect(f, h) {
			pkg.RespondNotOK(w, "badrequest")
			return
		}

		defer f.Close()
		filename := pkg.SaveFile(f, h.Filename)
		if filename != "" {
			post.ImageURL = "/pics/" + filename
		}
	}

	// check if the post content is correct
	err := pkg.ArePostInfosCorrect(post)
	if err != nil {
		pkg.RespondNotOK(w, "badrequest")
		return
	}
	// insert the post into the DB
	post, er = c.DB.InsertPostDB(userID, post)
	if er != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println(post.UserImageProfile)
	post.NumberOfComments = 0

	pkg.RespondOK(w, post, "post")

}
