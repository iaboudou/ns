package controllers

import (
	"fmt"

	handlers "rtf/controllers/handlers/chat"
	"rtf/models"
	"rtf/pkg/db/sqlite"
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

			case "mark_read":
				sqlite.MarkRead(db, msg.SenderID, msg.ReceiverID)

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

			case "load_group_history":
				GetOldGroupMessages(clients, db, msg)

			case "chat":
				Chat(clients, db, msg)

			case "group_chat":
				GroupChat(clients, db, msg)

			case "get_chat_users":
				users, _ := sqlite.GetChatUsers(db, msg.SenderID)

				cs, ok := clients[msg.SenderID]
				if !ok {
					continue
				}
				for _, conn := range cs {
					conn.Mu.Lock()
					conn.Ws.WriteJSON(map[string]any{
						"event": "chat_users",
						"users": users,
					})
					conn.Mu.Unlock()
				}

			case "get_notifications":
				GetNotificationsWS(clients, db, msg)

			case "get_unread_notifications_count":
				GetUnreadNotificationCountWS(clients, db, msg)
			}

		case client := <-hub.Disconnect:
			handlers.Disconnect(clients, client)

		case notif := <-hub.Notify:
			SendFollowNotification(clients, db, notif.FromUser, notif.ToUserID, notif.NotifType)
		}
	}
}
