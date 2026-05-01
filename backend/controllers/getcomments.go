package controllers

import (
	"encoding/json"
	"net/http"
	"rtf/help"
)

// this functionne to get all the comments
func (c *Controller) GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		help.RespondNotOK(w, "notallowed")
		return
	}

	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		help.RespondNotOK(w, "unauthorized")
		return
	}

	var req struct {
		Offset int    `json:"offset"`
		PostID string `json:"post_id"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	comments, er := c.DB.Get10PostComments(req.PostID, req.Offset)
	if er != nil {
		help.RespondNotOK(w, "server-error")
		return
	}

	if len(comments) == 0 {
		help.RespondOK(w, nil, "")
		return
	}

	help.RespondOK(w, comments, "comments")
}
