package main

import (
	"fmt"
	"log"
	"net/http"

	"rtf/controllers"
	"rtf/models"
	"rtf/pkg/db/sqlite"
	"rtf/routes"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, er := sqlite.InitDB()
	if er != nil {
		log.Fatal(er)
	}

	// initialize
	r := &sqlite.Repo{Db: database}

	hub := &models.Hub{
		Connect:    make(chan *models.Client),
		Disconnect: make(chan *models.Client),
		Broadcast:  make(chan models.Message),
		Notif:      make(chan models.Notification),
	}

	controller := &controllers.Controller{
		DB:  r,
		Hub: hub,
	}

	go controller.RunBroker()

	handler := &routes.Handler{
		Repo:    r,
		LastRL:  make(map[string]*routes.RateLimiter),
		Cntrlrs: controller,
	}

	mux := http.NewServeMux()
	routes.Routes(mux, handler)

	server := http.Server{
		Addr:    "0.0.0.0:4001",
		Handler: mux,
	}
	fmt.Println("http://localhost:4001")

	log.Fatal(server.ListenAndServe())
}
