package controllers

import (
	"net/http"

	handlers "rtf/controllers/handlers/groups"
	"rtf/help"
)

func (c *Controller) Groups(w http.ResponseWriter, r *http.Request) {
	segments := help.GetPathSegments(r.URL.Path)
	userID := r.Context().Value("userID").(string)

	switch r.Method {

	case http.MethodGet:
		c.handleGetGroups(w, r, segments, userID)

	case http.MethodPost:
		c.handlePostGroups(w, r, segments, userID)

	case http.MethodPatch:
		c.handlePatchGroups(w, r, segments, userID)

	case http.MethodDelete:
		c.handleDeleteGroups(w, segments, userID)

	default:
		help.RespondMethodNotAllowed(w)
	}
}

func (c *Controller) handleGetGroups(w http.ResponseWriter, r *http.Request, s []string, userID string) {
	switch {

	// GET /api/groups?want=...
	case len(s) == 0:
		switch r.URL.Query().Get("want") {
		case "mine":
			handlers.GetJoinedGroups(w, r, c.DB.Db)
		case "discover":
			handlers.GetUnknownGroups(w, r, c.DB.Db)
		case "invites":
			handlers.GetGroupInvites(w, r, c.DB.Db, userID)
		default:
			help.RespondNotFound(w, "This ressource doesn't exist")
		}

	// GET /api/groups/:groupID
	case len(s) == 1:
		handlers.GetGroupBasicInfo(w, r, c.DB.Db, s[0])

	// GET /api/groups/:groupID/requests
	case len(s) == 2 && s[1] == "requests":
		handlers.GetGroupRequests(w, r, c.DB.Db, s[0])

	// GET /api/groups/:groupID/events
	case len(s) == 2 && s[1] == "events":
		handlers.GetGroupEvents(w, r, c.DB.Db, s[0], userID)

	// GET /api/groups/:groupID/users
	case len(s) == 2 && s[1] == "users":
		handlers.GetUsers(w, r, c.DB.Db, s[0], userID)

	default:
		help.RespondNotFound(w, "This ressource doesn't exist")
	}
}

func (c *Controller) handlePostGroups(w http.ResponseWriter, r *http.Request, s []string, userID string) {
	switch {

	// POST /api/groups
	case len(s) == 0:
		handlers.CreateGroup(w, r, c.DB.Db, userID)

	// POST /api/groups/:groupID/requests
	case len(s) == 2 && s[1] == "requests":
		handlers.JoinGroup(w, r, c.Hub, c.DB.Db, s[0], userID, "request")

	// POST /api/groups/:groupID/events
	case len(s) == 2 && s[1] == "events":
		handlers.CreateEvent(w, r, c.DB.Db, c.Hub, s[0], userID)

	// POST /api/groups/:groupID/invites
	case len(s) == 2 && s[1] == "invites":
		handlers.InviteUser(w, r, c.Hub, c.DB.Db, s[0])

	default:
		help.RespondNotFound(w, "This ressource doesn't exist")
	}
}

func (c *Controller) handlePatchGroups(w http.ResponseWriter, r *http.Request, s []string, userID string) {
	switch {

	// PATCH /api/groups/:groupID/requests
	case len(s) == 2 && s[1] == "requests":
		handlers.ManageRequest(w, r, c.DB.Db, s[0], userID)

	// PATCH /api/groups/:groupID/events/:eventID
	case len(s) == 3 && s[1] == "events":
		handlers.VoteEvent(w, r, c.DB.Db, s[2], userID)

	default:
		help.RespondNotFound(w, "This ressource doesn't exist")
	}
}

func (c *Controller) handleDeleteGroups(w http.ResponseWriter, s []string, userID string) {
	switch {

	// DELETE /api/groups/:groupID
	case len(s) == 1:
		handlers.DeleteGroup(w, c.DB.Db, s[0], userID)

	// DELETE /api/groups/:groupID/me
	case len(s) == 2 && s[1] == "me":
		handlers.LeaveGroup(w, c.DB.Db, s[0], userID)

	default:
		help.RespondNotFound(w, "This ressource doesn't exist")
	}
}
