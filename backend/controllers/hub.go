package controllers

import (
	"fmt"

	handlers "rtf/controllers/handlers/chat"
	"rtf/models"
)

func (c *Controller) RunBroker() {
	clients := make(map[string][]*models.Client)
	db := c.DB.Db
	hub := c.Hub

	for {
		select {

		case client := <-hub.Connect:
			clients[client.ID] = append(clients[client.ID], client)

			err := handlers.Connect(clients, db, client)
			if err != nil {
				fmt.Println("broker: connect error:", err)
				continue
			}

		case msg := <-hub.Broadcast:
			switch msg.Type {
			// 	case "reconnect":
			// 		err := Reconnect(clients, db, msg.Sender)
			// 		if err != nil {
			// 			fmt.Printf("broker: failed to reconnect user: %v because of: %v\n", msg.Sender, err)
			// 		}

			// 	case "mark_read":
			// 		err := sqlite.MarkRead(db, msg.Sender, msg.Receiver)
			// 		if err != nil {
			// 			fmt.Println("broker: failed to mark read:", err)
			// 		}

			// 	case "get_unread":
			// 		err := GetUnread(clients, db, msg)
			// 		if err != nil {
			// 			fmt.Println("broker: failed to get unread messages:", err)
			// 		}

			case "load_history":
				err := GetOldMessages(clients, db, msg)
				if err != nil {
					fmt.Println("broker: failed to load history:", err)
				}

			case "chat":
				err := Chat(clients, db, msg)
				if err != nil {
					fmt.Println("broker: failed to send message :", err)
				}

				// 	case "typing":
				// 		Type(clients, msg.Receiver, msg.Sender)

				// 	case "stop-typing":
				// 		StopType(clients, msg.Receiver, msg.Sender)
			}

		case client := <-hub.Disconnect:
			handlers.Disconnect(clients, client)
		}
	}
}
