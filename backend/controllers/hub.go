package controllers

import (
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
				continue
			}

		case msg := <-hub.Broadcast:
			switch msg.Type {
			case "mark_read":
				sqlite.MarkRead(db, msg.SenderID, msg.ReceiverID)
				if cs, ok := clients[msg.SenderID]; ok {
					for _, conn := range cs {
						conn.Mu.Lock()
						conn.Ws.WriteJSON(map[string]any{
							"event": "messages_read",
							"data": map[string]any{
								"receiver_Id": msg.ReceiverID,
							},
						})
						conn.Mu.Unlock()
					}
				}

			case "load_history":
				GetOldMessages(clients, db, msg)

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
				GetNotificationsWS(clients, db, msg.SenderID)

			case "mark_notif_read":
				db.Exec(`UPDATE notification_users SET is_read = 1 WHERE notification_id = ? AND user_id = ?`, msg.ID, msg.SenderID)

			case "unread_notif":
				GetUnreadNotificationCountWS(clients, db, msg)
			}

		case notif := <-hub.Notif:

			switch notif.Type {
			case "event":
				BroadCastEventCreation(db, clients, &notif)
			case "group_invite":
				Notify(db, clients, &notif)
			case "group_request":
				Notify(db, clients, &notif)
			case "follow_request":
				Notify(db, clients, &notif)
			}
		case client := <-hub.Disconnect:
			handlers.Disconnect(clients, client)
		}
	}
}
