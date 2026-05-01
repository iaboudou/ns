package controllers

import (
	"net/http"
	"rtf/models"
	"rtf/pkg"
)

func (c *Controller) Getpersonalinfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		pkg.RespondNotOK(w, "unauthorized")
		return
	}

	id := r.URL.Query().Get("id")

	var user models.User
	var er error
	if id == userID {
		user, er = c.DB.GetPeronalInfoFromDB(userID)
		if er != nil {
			pkg.RespondNotOK(w, "server-error")
			return
		}
	} else {
		user, er = c.DB.GetPeronalInfoFromDB(id)
		if er != nil {
			pkg.RespondNotOK(w, "server-error")
			return
		}

		er = c.DB.IstheUserFreind(&user, userID)
		if er != nil {
			pkg.RespondNotOK(w, "badrequest")
			return
		}
	}

	pkg.RespondOK(w, user, "user")
}
