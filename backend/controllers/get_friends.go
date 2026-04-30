package controllers

import (
	"net/http"

	"rtf/pkg"
)

func (c *Controller) GetFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pkg.RespondNotOK(w, "notallowed")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		pkg.RespondNotOK(w, "unauthorized")
		return
	}

	users, er := c.DB.GetFriendsDB(userID)
	if er != nil {
		pkg.RespondNotOK(w, "server-error")
		return
	}

	pkg.RespondOK(w, users, "users")
}
