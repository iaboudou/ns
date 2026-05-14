package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"rtf/help"
	"rtf/models"
)

// get a list of posts from the DB and render it to the front
func (c *Controller) GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		help.RespondNotOK(w, "notallowed")
		return
	}

	viewerID, ok := r.Context().Value("userID").(string)
	if !ok || viewerID == "" {
		help.RespondNotOK(w, "unauthorized")
		return
	}

	q := r.URL.Query()
	offset, _ := strconv.Atoi(q.Get("offset"))
	page := q.Get("page")
	section := q.Get("section")
	reqUserID := q.Get("user_id")
	groupID := q.Get("group_id")

	if page == "profile-me-posts" || page == "profile-me-activity" {
		reqUserID = viewerID
	}

	if groupID != "" && groupID != "me" && page != "profille-other-posts" {
		var exists bool
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
		err = c.DB.Db.QueryRow(`SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?`, viewerID, groupID).Scan(&isMember)
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

	// viewer is the user want to see the posts
	// reqUserID is the user who owned the posts
	if page == "profille-other-posts" || reqUserID != viewerID {
		user := models.User{ID: reqUserID}
		if err := c.DB.IstheUserFreind(&user, viewerID); err != nil {
			help.RespondNotOK(w, "forbidden")
			return
		}
	}

	posts, err := c.DB.Get10PostsfromDB(page, section, viewerID, reqUserID, groupID, offset)
	if err != nil {
		help.RespondNotOK(w, "server-error")
		return
	}

	help.RespondOK(w, posts, "posts")
}
