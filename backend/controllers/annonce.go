package controllers

// import "rtf/help"

// func (c *Controller) Anounce() {
// 	c.Ws.Mu.Lock()
// 	if len(c.Ws.Clients) == 0 {
// 		c.Ws.Mu.Unlock()
// 		return
// 	}

// 	clients := make([]*UserWS, 0, len(c.Ws.Clients))
// 	for _, client := range c.Ws.Clients {
// 		clients = append(clients, client)
// 	}
// 	c.Ws.Mu.Unlock()

// 	for _, client := range clients {
// 		usersinfo, er := c.DB.GetUsersInfoFor(client.UserInfo.ID, true)
// 		if er != nil {
// 			continue
// 		}
// 		usersinfo = help.SortUsers(usersinfo)

// 		c.Ws.Mu.Lock()
// 		for i, u := range usersinfo {
// 			if _, ok := c.Ws.Clients[u.ID]; ok {
// 				usersinfo[i].IsOnline = true
// 			} else {
// 				usersinfo[i].IsOnline = false
// 			}
// 		}

// 		if client != nil && client.Chan != nil {
// 			select {
// 			case client.Chan <- map[string]interface{}{
// 				"type": "ws_users_info_for_user",
// 				"data": usersinfo,
// 			}:
// 			default:
// 			}
// 		}
// 		c.Ws.Mu.Unlock()
// 	}
// }
