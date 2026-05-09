package chat

import (
	"database/sql"

	"rtf/models"
)

func Connect(clients map[string][]*models.Client, db *sql.DB, newClient *models.Client) error {
	onlineUsers := []string{}

	for clientID := range clients {
		if clientID != newClient.ID {
			onlineUsers = append(onlineUsers, clientID)
		}
	}

	newClient.Mu.Lock()
	newClient.Ws.WriteJSON(map[string]any{
		"event": "online_users",
		"users": onlineUsers,
	})
	newClient.Mu.Unlock()

	// Only broadcast join when the user goes from offline to online.
	if len(clients[newClient.ID]) == 1 {
		for clientID, cs := range clients {
			if clientID == newClient.ID {
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
