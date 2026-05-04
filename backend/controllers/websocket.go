package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"rtf/help"
	"rtf/models"

	"github.com/gorilla/websocket"
)

func (c *Controller) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	user, er := c.DB.CheckSessionExistance(r)

	userID := user.ID

	if er != nil {
		help.RespondNotOK(w, "unauthorized")
		return
	}

	hub := c.Hub

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &models.Client{
		ID:          userID,
		Ws:          ws,
		Mu:          &sync.Mutex{},
		Description: &user,
	}

	hub.Connect <- client
	for {
		_, payload, err := ws.ReadMessage()
		fmt.Println("here: ", string(payload))
		if err != nil {
			client.Ws.Close()
			hub.Disconnect <- client
			return
		}

		var msg models.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			continue
		}

		msg.SenderID = client.ID
		hub.Broadcast <- msg
	}
}
