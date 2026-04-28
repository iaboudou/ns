package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"rtf/models"
	"rtf/pkg"
)

// this handler handles the  user registration. it expects a POST request, and it returns a JSON response with (error or success)
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		pkg.RespondNotOK(w, "badrequest")
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("here1")
		pkg.RespondNotOK(w, "badrequest")
		return
	}

	user := pkg.USERDATA(r)
	fmt.Println(user)

	file, header, err := r.FormFile("avatar")
	if err == nil && header.Size != 0 {
		defer file.Close()
		if !pkg.IsPictureFormatCorrect(file, header) {
			pkg.RespondNotOK(w, "badrequest")
			return
		}
		filename := pkg.SaveFile(file, header.Filename)
		if filename == "" {
			pkg.RespondNotOK(w, "badrequest")
			return
		}
		user.Avatar = filename
	}

	if err := pkg.CanInsertUser(&user); err != nil {
		pkg.RespondNotOK(w, err.Error())
		return
	}

	if err := c.DB.IsUserAlreadyExist(&user); err != nil {
		pkg.RespondNotOK(w, err.Error())
		return
	}

	if err := c.DB.InsertUserDB(user); err != nil {
		pkg.RespondNotOK(w, "server-error")
		return
	}

	pkg.RespondOK(w, nil, "")
	c.Anounce()
}

// this handler is for login. it expects a POST request, and it returns a JSON response with (error or success)
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid fields", http.StatusBadRequest)
		return
	}

	userID, er := c.DB.IsUserExist(&user)
	if er != nil {
		switch er.Error() {
		case "SERVER ERROR":
			http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
			return
		default:
			http.Error(w, fmt.Sprintf("%s", er.Error()), http.StatusBadRequest)
		}
		return
	}

	user, er = c.DB.GetUserInfos(userID)
	if er != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}

	a, er := c.DB.SetUserSession(w, userID)
	if er != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    a[0].(string),
		Path:     "/",
		HttpOnly: true,
		Expires:  a[1].(time.Time),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":    user,
		"success": "logged in successfully",
	})
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
		return
	}

	c.DB.DisconnectUser(userID)

	c.Ws.Clients[userID].RemoveUserWS(c.Ws, userID)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
		MaxAge:  -1,
	})
}
