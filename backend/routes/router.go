package routes

import (
	"net/http"
)

func Routes(mux *http.ServeMux, handler *Handler) {
	//
	routes := map[string]http.HandlerFunc{
		"/logout":        handler.Cntrlrs.Logout,
		"/createpost":    handler.Cntrlrs.CreatePost,
		"/createcomment": handler.Cntrlrs.CreateComment,
		"/getposts":      handler.Cntrlrs.GetPosts,
		"/getcomments":   handler.Cntrlrs.GetComments,
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
	}
	for path, h := range ws {
		mux.HandleFunc(path, handler.CORSMiddleware(h))
	}
}
