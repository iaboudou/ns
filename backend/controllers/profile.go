package controllers

import (
	"net/http"
	"rtf/pkg"
)

func (c *Controller) SwitchAccountPrivacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		pkg.RespondNotOK(w, "notallowed")
		return
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		pkg.RespondNotOK(w, "unauthorized")
		return
	}

	er := c.DB.SwitchAccountPrivacyinDB(userID)
	if er != nil {
		pkg.RespondNotOK(w, "server-error")
		return
	}

	pkg.RespondOK(w, nil, "")
}
