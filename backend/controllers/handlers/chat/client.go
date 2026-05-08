package chat

import (
	"database/sql"

	"rtf/models"
)

func Connect(clients map[string][]*models.Client, db *sql.DB, newClient *models.Client) error {
	onlineUsers := []string{}
	var count int

	for userID := range clients {
		onlineUsers = append(onlineUsers, userID)
		err := db.QueryRow(`SELECT COUNT(*) FROM notification_users WHERE user_id = ? AND is_read = 0`, userID).Scan(&count)
		if err != nil {
			return err
		}
	}

	for clientId, cs := range clients {
		if clientId == newClient.ID {
			if len(cs) == 0 {
				continue
			}
			cs[0].Mu.Lock()
			cs[0].Ws.WriteJSON(map[string]any{
				"event": "online_users",
				"users": onlineUsers,
				"count": count,
			})
			cs[0].Mu.Unlock()
		}
	}
	return nil
}

func Disconnect(clients map[string][]*models.Client, client *models.Client) {
	id := client.ID
	cs := clients[id]

	// Remove specific connection
	for i, c := range cs {
		if c == client {
			clients[id] = append(cs[:i], cs[i+1:]...)
			break
		}
	}

	// If no connections left for this user, broadcast "leave"
	if len(clients[id]) == 0 {
		delete(clients, id)
		for clientsId, cs := range clients {
			if clientsId == id {
				continue
			}

			for _, c := range cs {
				c.Mu.Lock()
				c.Ws.WriteJSON(map[string]any{
					"event": "leave",
					"left":  id,
				})
				c.Mu.Unlock()
			}
		}
	}
}

func GetUnreadNotificationCountWS(clients map[string][]*models.Client, db *sql.DB, msg models.Message) error {
	userID := msg.SenderID

	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM notification_users WHERE user_id = ? AND is_read = 0`, userID).Scan(&count)
	if err != nil {
		return err
	}

	if cs, ok := clients[userID]; ok {
		for _, c := range cs {
			c.Mu.Lock()
			c.Ws.WriteJSON(map[string]any{
				"event": "unread_notifications_count",
				"count": count,
			})
			c.Mu.Unlock()
		}
	}
	return nil
}
