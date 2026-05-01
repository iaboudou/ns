package controllers

import (
	"net/http"
	"rtf/help"
)

func (c *Controller) SwitchAccountPrivacy(w http.ResponseWriter, r *http.Request) {
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

	er := c.DB.SwitchAccountPrivacyinDB(userID)
	if er != nil {
		help.RespondNotOK(w, "server-error")
		return
	}

	help.RespondOK(w, nil, "")
}
