package controllers

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"rtf/models"
	"rtf/pkg/db/sqlite"
)

func GetOldMessages(clients map[string][]*models.Client, db *sql.DB, msg models.Message) error {
	cs, ok := clients[msg.SenderID]
	if !ok {
		return nil
	}

	// mark everything read for this conversation before sending history
	err := sqlite.MarkRead(db, msg.SenderID, msg.ReceiverID)
	if err != nil {
		return err
	}

	messages, err := sqlite.SelectOldMessages(db, &msg)
	if err != nil {
		return err
	}

	for _, client := range cs {
		client.Mu.Lock()
		client.Ws.WriteJSON(map[string]any{
			"event":       "history",
			"messages":    messages,
			"hasMore":     len(messages) == 10,
			"receiver_id": msg.ReceiverID,
			"portKey":     msg.PortKey,
		})

		client.Mu.Unlock()
	}

	return nil
}


func Chat(clients map[string][]*models.Client, db *sql.DB, msg models.Message) error {
	if len(strings.TrimSpace(msg.Content)) == 0 {
		return errors.New("message is empty")
	}

	if len(msg.Content) > 2000 {
		return errors.New("message is too long")
	}

	msg.CreatedAt = time.Now().UnixMilli()

	err := sqlite.InsertNewMessage(db, &msg)
	if err != nil {
		return err
	}

	senderConns := clients[msg.SenderID]

	sender := senderConns[0].Description

	receiverConns, receiverOnline := clients[msg.ReceiverID]

	var receiver *models.User

	if receiverOnline {
		receiver = receiverConns[0].Description
	} else {
		u, err := sqlite.GetUserByID(db, msg.ReceiverID)
		if err != nil {
			return err
		}

		receiver = &u
	}

	payload := map[string]any{
		"event": "own_message",
		"message": map[string]any{
			"id":                msg.ID,
			"content":           msg.Content,
			"created_at":        msg.CreatedAt,
			"sender_id":         sender.ID,
			"receiver_id":       receiver.ID,
			"sender_nickname":   sender.Nickname,
			"receiver_nickname": receiver.Nickname,
			"sender_fullname":   sender.Firstname + " " + sender.Lastname,
			"receiver_fullname": receiver.Firstname + " " + receiver.Lastname,
			"sender_profile":    sender.ProfileImage,
			"receiver_profile":  receiver.ProfileImage,
		},
	}

	for _, conn := range senderConns {
		conn.Mu.Lock()
		conn.Ws.WriteJSON(payload)
		conn.Mu.Unlock()
	}

	if receiverOnline {
		payload["event"] = "other_message"
		for _, conn := range receiverConns {
			conn.Mu.Lock()
			conn.Ws.WriteJSON(payload)
			conn.Mu.Unlock()
		}
	}

	return nil
}


func GetOldGroupMessages(clients map[string][]*models.Client, db *sql.DB, msg models.Message) error {
	cs, ok := clients[msg.SenderID]
	if !ok {
		return nil
	}

	messages, err := sqlite.SelectOldGroupMessages(db, msg.ReceiverID, msg.LastReadTime)
	if err != nil {
		return err
	}

	for _, client := range cs {
		client.Mu.Lock()
		client.Ws.WriteJSON(map[string]any{
			"event":       "history",
			"messages":    messages,
			"hasMore":     len(messages) == 10,
			"receiver_Id": msg.ReceiverID,
			"isGroup":     true,
			"portKey":     msg.PortKey,
		})
		client.Mu.Unlock()
	}

	return nil
}

func GroupChat(clients map[string][]*models.Client, db *sql.DB, msg models.Message) error {
	if len(strings.TrimSpace(msg.Content)) == 0 {
		return errors.New("message is empty")
	}

	msg.CreatedAt = time.Now().UnixMilli()

	// insert message
	msgID, err := sqlite.InsertGroupMessage(db, msg.ReceiverID, msg.SenderID, msg.Content, msg.CreatedAt)
	if err != nil {
		return err
	}

	// get sender info
	senderUser, err := sqlite.GetUserByID(db, msg.SenderID)
	if err != nil {
		return err
	}

	// get group members
	members, err := sqlite.GetGroupMembers(db, msg.ReceiverID)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"event": "new_group_message",
		"message": map[string]any{
			"id":              msgID,
			"group_id":        msg.ReceiverID,
			"content":         msg.Content,
			"created_at":      msg.CreatedAt,
			"sender_id":       msg.SenderID,
			"sender_nickname": senderUser.Nickname,
			"sender_fullname": senderUser.Firstname + " " + senderUser.Lastname,
			"sender_profile":  senderUser.ProfileImage,
			"is_group":        true,
		},
	}

	// broadcast to all online members
	for _, memberID := range members {
		if memberConns, online := clients[memberID]; online {
			for _, conn := range memberConns {
				conn.Mu.Lock()
				conn.Ws.WriteJSON(payload)
				conn.Mu.Unlock()
			}
		}
	}

	return nil
}
