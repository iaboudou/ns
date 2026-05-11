package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"rtf/help"
	"rtf/models"
)

// this handler handles the  user registration. it expects a POST request, and it returns a JSON response with (error or success)
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		help.Respond(w, &models.Response{
			Code:    http.StatusMethodNotAllowed,
			Message: "not allowed",
		})
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid file image",
		})
		return
	}

	user := help.USERDATA(r)

	file, header, err := r.FormFile("avatar")
	if err == nil && header.Size != 0 {
		defer file.Close()
		if !help.IsPictureFormatCorrect(file, header) {
			help.Respond(w, &models.Response{
				Code:    http.StatusBadRequest,
				Message: "invalid image",
			})
			return
		}
		filename := help.SaveFile(file, header.Filename)
		if filename == "" {
			help.Respond(w, &models.Response{
				Code:    http.StatusBadRequest,
				Message: "invalid credential",
			})
			return
		}
		user.Avatar = filename
	}

	if err := help.CanInsertUser(&user); err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid credential",
		})
		return
	}

	if err := c.DB.IsUserAlreadyExist(&user); err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "credential already used",
		})
		return
	}

	if err := c.DB.InsertUserDB(user); err != nil {
		help.Respond(w, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "server-error",
		})
		return
	}

	help.RespondOK(w, nil, "")
}

// this handler is for login. it expects a POST request, and it returns a JSON response with (error or success)
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		help.RespondNotOK(w, "notallowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		help.RespondNotOK(w, "badrequest")

		help.Respond(w, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid credential",
		})
		return
	}

	userID, er := c.DB.IsUserExist(req.Email, req.Password)
	if er != nil {
		help.RespondNotOK(w, "badrequest")

		help.Respond(w, &models.Response{
			Code: http.StatusBadRequest,

			Message: "invalid credential",
		})
		return
	}

	user, er := c.DB.GetUserInfos(userID)
	if er != nil {
		help.RespondNotOK(w, "server-error")
		return
	}

	sessionID, expiresAt, er := c.DB.SetUserSession(w, userID)
	if er != nil {
		help.RespondNotOK(w, "server-error")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
	help.RespondOK(w, user, "user")
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		help.Respond(w, &models.Response{
			Code:    http.StatusMethodNotAllowed,
			Message: "not allowed",
		})
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		c.DB.DisconnectUser(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
		MaxAge:  -1,
	})

	help.RespondOK(w, nil, "")
}