package controllers

import (
	"rtf/models"
	"rtf/pkg/db/sqlite"
)

type Controller struct {
	DB  *sqlite.Repo
	Hub *models.Hub
}
