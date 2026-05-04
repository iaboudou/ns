package help

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"strings"

	"rtf/models"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

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
	if len(comment.ImageURL) > 0 && len(comment.Content) == 0 {
		return true
	}
	return len(comment.Content) <= 500 && len(comment.Content) > 0
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
