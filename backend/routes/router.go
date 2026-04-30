package routes

import (
	"net/http"
)

func Routes(mux *http.ServeMux, handler *Handler) {
	//
	routes := map[string]http.HandlerFunc{
		"/api/logout":                 handler.Cntrlrs.Logout,
		"/api/createpost":             handler.Cntrlrs.CreatePost,
		"/api/createcomment":          handler.Cntrlrs.CreateComment,
		"/api/getposts":               handler.Cntrlrs.GetPosts,
		"/api/getcomments":            handler.Cntrlrs.GetComments,
		"/api/follow":                 handler.Cntrlrs.Follow,
		"/api/getfriends":             handler.Cntrlrs.GetFriends,
		"/api/getpersonalinfo":        handler.Cntrlrs.Getpersonalinfo,
		"/api/switchaccountprivacy":   handler.Cntrlrs.SwitchAccountPrivacy,
		"/api/getsuggestionfollowers": handler.Cntrlrs.GetSuggestionFollowers,
	}

	for path, h := range routes {
		mux.HandleFunc(path, handler.CORSMiddleware(handler.RateLimit(handler.Middleware(h))))
	}

	// home page, login and register routes
	LRroutes := map[string]http.HandlerFunc{
		"/api/login":    handler.Cntrlrs.Login,
		"/api/register": handler.Cntrlrs.Register,
	}
	for path, h := range LRroutes {
		mux.HandleFunc(path, handler.CORSMiddleware(handler.RateLimit(h)))
	}

	//
	ws := map[string]http.HandlerFunc{
		"/ws":         handler.Cntrlrs.WebSocket,
		"/hassession": handler.Cntrlrs.HasSession,
		"/pics/":      handler.Cntrlrs.ServePictures,
	}
	for path, h := range ws {
		mux.HandleFunc(path, handler.CORSMiddleware(h))
	}
}
