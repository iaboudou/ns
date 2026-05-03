package chat

import (
	"database/sql"

	"rtf/models"
)

func Connect(clients map[string][]*models.Client, db *sql.DB, newClient *models.Client) error {
	onlineUsers := []string{}
	for userID := range clients {
		onlineUsers = append(onlineUsers, userID)
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
			})
			cs[0].Mu.Unlock()
			continue
		}

		for _, c := range cs {
			c.Mu.Lock()
			c.Ws.WriteJSON(map[string]any{
				"event":    "join",
				"newcomer": newClient.ID,
			})
			c.Mu.Unlock()
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
