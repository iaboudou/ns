package controllers

import (
	"net/http"

	"rtf/help"
	"rtf/models"
)

func (c *Controller) UpdateLastNotif(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		help.RespondMethodNotAllowed(w)
		return
	}

	userID := r.Context().Value("userID").(string)

	_, err := c.DB.Db.Exec("UPDATE users SET last_notif_seen = CURRENT_TIMESTAMP WHERE id = ?", userID)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
	})
}
