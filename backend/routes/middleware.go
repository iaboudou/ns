package routes

import (
	"context"
	"net/http"
	"rtf/help"
)

// return StatusUnauthorized if the user not loggedin
func (h *Handler) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// check session existance
		user, err := h.Repo.CheckSessionExistance(r)
		if err != nil {
			help.RespondNotOK(w, "unauthorized")
			return
		}
		if len(user.ID) == 0 {
			help.RespondNotOK(w, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), "userID", user.ID)
		next(w, r.WithContext(ctx))
	}
}

func (h *Handler) CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
