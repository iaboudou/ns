package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"rtf/help"
	"rtf/models"
)

func (c *Controller) Follow(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		demand := r.URL.Query().Get("want")
		targetID := r.URL.Query().Get("id")

		if targetID == "" {
			help.Respond(w, &models.Response{
				Code:    http.StatusBadRequest,
				Message: "missing user id",
			})
			return
		}

		var data any
		var err error
		q := r.URL.Query().Get("q")

		switch demand {
		case "followers":
			data, err = c.DB.GetFollowersDB(targetID, q)
		case "following":
			data, err = c.DB.GetFollowingDB(targetID, q)
		default:
			help.Respond(w, &models.Response{
				Code:    http.StatusBadRequest,
				Message: "unknown demand",
			})
			return
		}

		if err != nil {
			help.RespondNotOK(w, "server-error")
			return
		}

		help.Respond(w, &models.Response{
			Code: http.StatusOK,
			Data: data,
		})
		return
	}

	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body map[string]any
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
		return
	}

	message, err := c.DB.FollowUserDB(userID, body["followed_id"].(string))
	if err != nil {
		help.RespondServerError(w)
		return
	}

	// Send notification if this was a new follow
	if message == "follow have been successfully" || message == "request have been sent" {
		senderUser, err := c.DB.GetUserByIDDB(userID)
		if err == nil {
			notifType := "follow_accepted"
			if message == "request have been sent" {
				notifType = "follow_request"
			}
			go func() {
				c.Hub.Notify <- models.FollowNotif{
					FromUser:  senderUser,
					ToUserID:  body["followed_id"].(string),
					NotifType: notifType,
				}
			}()
		}
	}

	help.Respond(w, &models.Response{
		Code:    http.StatusOK,
		Message: message,
	})
}

func (c *Controller) GetSuggestionFollowers(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		help.Respond(w, &models.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	suggestions, err := c.DB.GetSuggestionUsersDB(userID)
	if err != nil {
		fmt.Printf("Error getting suggestions: %v\n", err)
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: suggestions,
	})
}

func (c *Controller) ManageFollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		help.Respond(w, &models.Response{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed",
		})
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		help.Respond(w, &models.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	var body struct {
		FollowerID string `json:"follower_id"`
		Decision   string `json:"decision"` // "accepted" or "rejected"
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
		return
	}

	err = c.DB.ManageFollowDB(body.FollowerID, userID, body.Decision)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code:    http.StatusOK,
		Message: "request successfully " + body.Decision,
	})
}

func (c *Controller) GetFollowRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		help.Respond(w, &models.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	requests, err := c.DB.GetFollowRequestsDB(userID)
	if err != nil {
		help.RespondServerError(w)
		return
	}

	help.Respond(w, &models.Response{
		Code: http.StatusOK,
		Data: requests,
	})
}
