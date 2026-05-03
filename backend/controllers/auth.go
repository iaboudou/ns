package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"rtf/help"
)

// this handler handles the  user registration. it expects a POST request, and it returns a JSON response with (error or success)
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		help.RespondNotOK(w, "badrequest")
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		help.RespondNotOK(w, "badrequest")
		return
	}

	user := help.USERDATA(r)

	file, header, err := r.FormFile("avatar")
	if err == nil && header.Size != 0 {
		defer file.Close()
		if !help.IsPictureFormatCorrect(file, header) {
			help.RespondNotOK(w, "badrequest")
			return
		}
		filename := help.SaveFile(file, header.Filename)
		if filename == "" {
			help.RespondNotOK(w, "badrequest")
			return
		}
		user.Avatar = filename
	}

	if err := help.CanInsertUser(&user); err != nil {
		help.RespondNotOK(w, err.Error())
		return
	}

	if err := c.DB.IsUserAlreadyExist(&user); err != nil {
		help.RespondNotOK(w, err.Error())
		return
	}

	if err := c.DB.InsertUserDB(user); err != nil {
		help.RespondNotOK(w, "server-error")
		return
	}

	help.RespondOK(w, nil, "")
	// c.Anounce()
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
		return
	}

	userID, er := c.DB.IsUserExist(req.Email, req.Password)
	if er != nil {
		help.RespondNotOK(w, "badrequest")
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

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		help.RespondNotOK(w, "unauthorized")
		return
	}

	c.DB.DisconnectUser(userID)

	// if ws, ok := c.Ws.Clients[userID]; ok && ws != nil {
	// 	ws.RemoveUserWS(c.Ws, userID)
	// }

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
		MaxAge:  -1,
	})

	help.RespondOK(w, nil, "")
}
