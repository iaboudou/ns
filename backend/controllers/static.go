package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"rtf/help"
	"strings"
)

func (hand *Controller) ServePictures(w http.ResponseWriter, r *http.Request) {
	f := strings.TrimPrefix(r.URL.Path, "/pics/")
	path := filepath.Join("./db/pics", f)
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		help.RespondNotOK(w, "notfound")
		return
	}

	http.ServeFile(w, r, path)
}
