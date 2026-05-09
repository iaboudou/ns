package controllers

import (
	"net/http"

	"rtf/help"
	"rtf/models"
)

func (c *Controller) GetFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		help.RespondNotOK(w, "notallowed")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		help.RespondNotOK(w, "unauthorized")
		return
	}

	q := r.URL.Query().Get("q")
	isallusers := r.URL.Query().Get("allusers")

	var users []models.FollowSuggestion
	var er error
	if isallusers == "true" {
		users, er = c.DB.GetUsers(userID, q)
	}else {
		users, er = c.DB.GetFollowersDB(userID, q)
	}

	if er != nil {
		help.RespondNotOK(w, "server-error")
		return
	}

	help.RespondOK(w, users, "users")
}
