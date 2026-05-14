package handlers

import (
	"database/sql"
	"net/http"

	"rtf/help"
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func GetUnknownGroups(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.Context().Value("userID").(string)
	filter := r.URL.Query().Get("search")
	last := r.URL.Query().Get("last")
	lastID := r.URL.Query().Get("lastId")

	groups, err := sqlite.SelectUnknownGroupsInDB(db, userID, filter, last, lastID)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: groups,
	})
}
