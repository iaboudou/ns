package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"rtf/config"
)

// / USER
type U struct {
	ID        string
	Nickname  string
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Birthday  string
	Gender    string
	About     string
	Avatar    string
	SessionID string
	AccountPrivacy  bool
}

func USERDATA(r *http.Request) U {
	return U{
		Nickname:  r.FormValue("nickname"),
		Firstname: r.FormValue("firstname"),
		Lastname:  r.FormValue("lastname"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		Birthday:  r.FormValue("dob"),
		Gender:    r.FormValue("gender"),
		About:     r.FormValue("about"),
	}
}

func CanInsertUser(user *U) error {
	//
	user.Email = strings.TrimSpace(user.Email)
	user.Firstname = strings.TrimSpace(user.Firstname)
	user.Lastname = strings.TrimSpace(user.Lastname)
	user.Nickname = strings.TrimSpace(user.Nickname)
	user.Birthday = strings.TrimSpace(user.Birthday)
	user.Gender = strings.TrimSpace(user.Gender)
	user.About = strings.TrimSpace(user.About)
	user.Avatar = strings.TrimSpace(user.Avatar)

	if len(user.Email) == 0 || len(user.Email) > config.MaxEmailLength {
		return fmt.Errorf("Email required, max %d chars", config.MaxEmailLength)
	}
	if len(user.Firstname) == 0 || len(user.Firstname) > config.MaxFirstnameLength || len(user.Firstname) < 2 {
		return fmt.Errorf("Firstname required, max %d chars", config.MaxFirstnameLength)
	}
	if len(user.Lastname) == 0 || len(user.Lastname) > config.MaxLastnameLength || len(user.Lastname) < 2 {
		return fmt.Errorf("Lastname required, max %d chars", config.MaxLastnameLength)
	}

	if len(user.Password) < config.MinPasswordLength || len(user.Password) > config.MaxPasswordLength {
		return fmt.Errorf("Password %d-%d chars", config.MinPasswordLength, config.MaxPasswordLength)
	}
	if user.Birthday == "" {
		return errors.New("Birthday required")
	}

	g := strings.ToLower(user.Gender)
	if g != "male" && g != "female" {
		return errors.New("Gender must be 'male' or 'female'")
	}
	user.Gender = g

	return nil
}
