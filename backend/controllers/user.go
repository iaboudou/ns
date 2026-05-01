package controllers

import (
	"net/http"
	"rtf/models"
	"rtf/help"
)

func (c *Controller) Getpersonalinfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		help.RespondNotOK(w, "unauthorized")
		return
	}

	id := r.URL.Query().Get("id")

	var user models.User
	var er error
	if id == userID {
		user, er = c.DB.GetPeronalInfoFromDB(userID)
		if er != nil {
			help.RespondNotOK(w, "server-error")
			return
		}
	} else {
		user, er = c.DB.GetPeronalInfoFromDB(id)
		if er != nil {
			help.RespondNotOK(w, "server-error")
			return
		}

		er = c.DB.IstheUserFreind(&user, userID)
		if er != nil {
			help.RespondNotOK(w, "badrequest")
			return
		}
	}

	help.RespondOK(w, user, "user")
}
