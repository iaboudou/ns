package controllers

import (
	"net/http"
	"rtf/models"
	"rtf/pkg"
	"strconv"
)

// get a list of posts from the DB and render it to the front
func (c *Controller) GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pkg.RespondNotOK(w, "notallowed")
		return
	}

	viewerID, ok := r.Context().Value("userID").(string)
	if !ok || viewerID == "" {
		pkg.RespondNotOK(w, "unauthorized")
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

	if page == "profille-other-posts" || reqUserID != viewerID {
		user := models.User{ID: reqUserID}
		if err := c.DB.IstheUserFreind(&user, viewerID); err != nil {
			pkg.RespondNotOK(w, "forbidden")
			return
		}
	}

	posts, err := c.DB.Get10PostsfromDB(page, section, viewerID, reqUserID, groupID, offset)
	if err != nil {
		pkg.RespondNotOK(w, "server-error")
		return
	}
	pkg.RespondOK(w, posts, "posts")
}
