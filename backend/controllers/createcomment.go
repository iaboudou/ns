package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

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
	path := r.FormValue("path")
	postID := r.FormValue("post_id")
	comment.PostID = postID

	// in case if want to create comment from /groups
	if strings.Contains(path, "groups/") {
		groupID := strings.Split(path, "/")[2]

		var exists int
		err := c.DB.Db.QueryRow(`SELECT 1 FROM groups WHERE id = ?`, groupID).Scan(&exists) // check group existence
		if err == sql.ErrNoRows {
			help.RespondNotFound(w, "This group doesn't exist")
			return
		}

		if err != nil {
			help.RespondServerError(w)
			return
		}

		var isMember bool
		err = c.DB.Db.QueryRow(`SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?`, userID, groupID).Scan(&isMember)
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
	}

	// if want to create comment from / or /profile

	var allowed_users, privacy string
	err := c.DB.Db.QueryRow(`SELECT privacy, allowed_users FROM posts WHERE id = ?`, comment.PostID).Scan(&privacy, &allowed_users)
	if err != nil {
		fmt.Println("err: ", err)
	}

	var creatorID string
	err = c.DB.Db.QueryRow(`SELECT user_id FROM posts WHERE id = ?`, postID).Scan(&creatorID)
	if err != nil {
		fmt.Println("er1: ", err)
		return
	}
	
	if privacy == "private" && creatorID != userID {
		if !strings.Contains(allowed_users, userID) {
			help.RespondNotOK(w, "badrequest")
			return
		}
	} else if privacy == "followers" && creatorID != userID {

		var follower int
		fmt.Println(userID, creatorID)
		err = c.DB.Db.QueryRow(`SELECT 1 FROM followers WHERE following_id = ? AND follower_id = ? `, userID, creatorID).Scan(&follower)
		if err != nil {
			fmt.Println(err)
			return
		}
		if follower <= 0 {
			return
		}
	}

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
