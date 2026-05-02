package controllers

import (
	"net/http"

	"rtf/help"
)

// this functione to check if the user has session
func (c *Controller) HasSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		help.RespondNotOK(w, "notallowed")
		return
	}

	_, er := c.DB.CheckSessionExistance(r)
	if er != nil {
		help.RespondNotOK(w, "unauthorized")
		return
	}
	help.RespondOK(w, nil, "")
}
