package controllers

import (
	"net/http"
	"rtf/pkg"
)

// this functione to check if the user has session
func (c *Controller) HasSession(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		pkg.RespondNotOK(w, "notallowed")
		return
	}

	_, er := c.DB.CheckSessionExistance(r)
	if er != nil {
		pkg.RespondNotOK(w, "unauthorized")
		return
	}

	pkg.RespondOK(w, nil, "")
}
