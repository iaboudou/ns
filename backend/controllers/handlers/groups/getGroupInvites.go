package handlers

import (
	"database/sql"
	"net/http"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func GetGroupInvites(w http.ResponseWriter, r *http.Request, db *sql.DB, userID string) {
	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")

	groups, err := sqlite.SelectGroupInvitesInDB(db, userID, last, lastID)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: 200,
		Data: groups,
	})
}
