package controllers

import (
	"encoding/json"
	"net/http"
	"rtf/pkg"
)

// this functionne to get all the comments
func (c *Controller) GetComments(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		pkg.RespondNotOK(w, "unauthorized")
		return
	}

	var req struct {
		Offset int    `json:"offset"`
		PostID string `json:"post_id"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	comments, er := c.DB.Get10PostComments(req.PostID, req.Offset)
	if er != nil {
		pkg.RespondNotOK(w, "server-error")
		return
	}

	if len(comments) == 0 {
		pkg.RespondOK(w, nil, "")
		return
	}

	pkg.RespondOK(w, comments, "comments")
}
