package pkg

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"rtf/config"
	"rtf/models"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// sort the users that we have a conversation with by the last message date and the rest of the users by their nickname
func SortUsers(usersinfo []models.UserInfo) []models.UserInfo {
	res := []models.UserInfo{}

	conv := []models.UserInfo{}
	vide := []models.UserInfo{}

	for _, i := range usersinfo {
		if i.LastMessage.ID >= 1 {
			conv = append(conv, i)
		} else {
			vide = append(vide, i)
		}
	}

	sort.Slice(conv, func(i, j int) bool {
		t1, _ := time.Parse("2006-01-02 15:04:05", conv[i].LastMessage.CreatedAt)
		t2, _ := time.Parse("2006-01-02 15:04:05", conv[j].LastMessage.CreatedAt)
		return t1.After(t2)
	})

	sort.Slice(vide, func(i, j int) bool {
		return vide[i].Nickname < vide[j].Nickname
	})

	res = append(res, conv...)
	res = append(res, vide...)
	return res
}

// this funcion check if the informations mutch the expected , any error found will be returned
func AreUserInfosCorret(user models.User) error {
	// empty feild
	if len(user.Nickname) == 0 ||
		len(user.Birthday) == 0 ||
		len(user.Gender) == 0 ||
		len(user.Firstname) == 0 ||
		len(user.Lastname) == 0 ||
		len(user.Email) == 0 ||
		len(user.Password) == 0 {
		return errors.New("all feilds are required")
	}

	// if user too young
	b, err := time.Parse("2006-01-02", user.Birthday)
	if err != nil {
		return errors.New("invalid date format")
	}
	now := time.Now().Unix()
	max := now - int64(60*60*24*365.25*200)
	legal := now - int64(60*60*24*365.25*15)
	birth_ms := b.Unix()

	if birth_ms > legal || birth_ms < max {
		return errors.New("you're not allowed to use this website")
	}

	// gender
	if user.Gender != "Male" &&
		user.Gender != "Female" &&
		user.Gender != "Other" {
		return errors.New("invalid gender")
	}

	// check the format of the firstname/lastname/nickname
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(user.Nickname) {
		return errors.New("invalid nickname format")
	}
	if !regexp.MustCompile(`^[a-zA-Z_]+$`).MatchString(user.Firstname) {
		return errors.New("invalid firstname format")
	}
	if !regexp.MustCompile(`^[a-zA-Z_]+$`).MatchString(user.Lastname) {
		return errors.New("invalid lastname format")
	}

	// email format
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	// more than max length
	if len(user.Nickname) > 20 || len(user.Firstname) > 20 || len(user.Lastname) > 20 ||
		len(user.Email) > 60 || len(user.Password) > 60 {
		return errors.New("feild too large")
	}

	if len(user.Nickname) < 4 || len(user.Firstname) < 4 || len(user.Lastname) < 4 ||
		len(user.Email) < 8 || len(user.Password) < 8 {
		return errors.New("feild too small")
	}

	return nil
}

// this function try hash the password with bcrypt , any error found will be returned
func HashPassword(password string) (string, error) {
	hashed, er := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if er != nil {
		return "", errors.New("SERVER ERROR")
	}
	return string(hashed), nil
}

// check if the post content is valid
func ArePostInfosCorrect(post models.Post) error {
	fmt.Println(post.Content, post.ImageURL, post.Privacy)
	if len(post.Content) > 600 {
		return errors.New("post too large")
	}
	if len(post.Content) == 0 && len(post.ImageURL) == 0 {
		return errors.New("post empty")
	}
	return nil
}

// check if the comment data is correct
func IsvalidComment(comment models.Comment) bool {
	return len(comment.Content) != 0 && len(comment.Content) < 500
}

// this function handle the rate limit for the messages
func MessageRLExceeded(count int, last time.Time) bool {
	if time.Since(last) > config.FiveSec {
		return false
	}
	return count >= config.Max_Messages
}

// check if the message format is correct
func TheMessageFormatIsCorrect(data map[string]interface{}) bool {
	content, ok := data["Content"].(string)
	if !ok {
		return false
	}
	receiverID, ok := data["ReceiverID"].(string)
	if !ok {
		return false
	}
	return len(content) > 0 && len(content) < 400 && len(receiverID) > 0 && len(receiverID) < 100
}

// this is a helper function return HTTP errors
func StatusError(w http.ResponseWriter, er error) {
	if er != nil {
		if er.Error() == "SERVER ERROR" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, er)))
		return
	}
}

// this function save a file in the pics folder and return its new name that was generated randomly
func SaveFile(r io.Reader, originalName string) string {
	randomID, err := uuid.NewV4()
	if err != nil {
		return ""
	}

	name := randomID.String() + filepath.Ext(originalName)
	cwd, _ := os.Getwd()
	fp := filepath.Join(cwd, "db", "pics", name)

	//
	os.MkdirAll("db/pics", os.ModePerm)

	out, err := os.Create(fp)
	if err != nil {
		return ""
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		return ""
	}

	return name
}

func IsPictureFormatCorrect(file multipart.File, header *multipart.FileHeader) bool {
	// check taille max 5MB
	if header.Size > 5<<20 {
		return false
	}

	//
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return false
	}
	fileType := http.DetectContentType(buf)
	if !allowedTypes[fileType] {
		return false
	}

	// reset read pointer
	file.Seek(0, io.SeekStart)

	// check extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExt[ext] {
		return false
	}

	return true
}
