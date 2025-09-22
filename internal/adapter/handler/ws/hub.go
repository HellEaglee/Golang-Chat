package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/google/uuid"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client

	BroadcastToRoom chan *RoomMessage

	ChatService    port.ChatService
	MessageService port.MessageService
}

func NewHub(cs port.ChatService, ms port.MessageService) *Hub {
	return &Hub{
		Clients:         make(map[*Client]bool),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		BroadcastToRoom: make(chan *RoomMessage),
		ChatService:     cs,
		MessageService:  ms,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			h.sendChatList(client)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case roomMsg := <-h.BroadcastToRoom:
			for client := range h.Clients {
				if client.InRoom(roomMsg.RoomID) {
					select {
					case client.Send <- roomMsg.Data:
					default:
						log.Printf("Client %s send buffer full — disconnecting", client.UserID)
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}

func (h *Hub) sendChatList(client *Client) {
	ctx := context.Background()
	chats, err := h.ChatService.GetChatsByUserID(ctx, client.UserID)
	if err != nil {
		log.Printf("Error fetching user chats: %v", err)
		return
	}

	chatList := make([]ChatListItem, len(chats))
	for i, c := range chats {
		participantsIds := make([]uuid.UUID, len(c.Participants))
		for j, p := range c.Participants {
			participantsIds[j] = p.UserID
		}

		chatList[i] = ChatListItem{
			ID:             c.ID,
			Name:           c.Name,
			IsGroup:        c.IsGroup,
			LastMessage:    c.LastMessage,
			LastMessageAt:  c.LastMessageAt,
			ParticipantIDs: participantsIds,
		}
	}

	data, _ := json.Marshal(Message{
		Type:    "chat_list",
		Payload: chatList,
	})

	client.Send <- data
}

func (h *Hub) HandleMessage(client *Client, msg Message) {
	ctx := context.Background()
	switch msg.Type {
	case "join":
		if roomIDStr, ok := msg.Payload.(string); ok {
			client.JoinRoom(roomIDStr)

			messages, err := h.MessageService.GetMessagesByChatID(ctx, roomIDStr)
			if err != nil {
				log.Printf("Error loading messages: %v", err)
				return
			}

			messageItems := make([]MessageItem, len(messages))
			for i, m := range messages {
				username := "Unknown"
				if m.User.Name != "" {
					username = m.User.Name
				}

				var replies []MessageItem
				for _, r := range m.Replies {
					replyUsername := "Unknown"
					if r.User.Name != "" {
						replyUsername = r.User.Name
					}
					replies = append(replies, MessageItem{
						ID:               r.ID,
						ChatID:           r.ChatID,
						UserID:           r.UserID,
						Username:         replyUsername,
						Text:             r.Text,
						IsEdited:         r.IsEdited,
						ReplyToMessageID: r.ReplyToMessageID,
						CreatedAt:        r.CreatedAt,
					})
				}

				messageItems[i] = MessageItem{
					ID:               m.ID,
					ChatID:           m.ChatID,
					UserID:           m.UserID,
					Username:         username,
					Text:             m.Text,
					IsEdited:         m.IsEdited,
					ReplyToMessageID: m.ReplyToMessageID,
					CreatedAt:        m.CreatedAt,
					Replies:          replies,
				}
			}

			historyData, _ := json.Marshal(Message{
				Type:    "message_history",
				Payload: messageItems,
			})

			client.Send <- historyData
		}

	case "leave":
		if roomIDStr, ok := msg.Payload.(string); ok {
			client.LeaveRoom(roomIDStr)
		}

	case "message":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			roomIDStr := payload["room_id"].(string)
			roomID, err := uuid.Parse(roomIDStr)
			if err != nil {
				return
			}

			text := payload["text"].(string)

			newMsg := domain.Message{
				ChatID:    roomID,
				UserID:    uuid.MustParse(client.UserID),
				Text:      text,
				IsEdited:  false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			savedMsg, err := h.MessageService.CreateMessage(ctx, &newMsg)
			if err != nil {
				log.Printf("Error saving message: %v", err)
				return
			}

			_, err = h.ChatService.UpdateLastMessage(ctx, roomIDStr, savedMsg.Text)
			if err != nil {
				log.Printf("Error updating last message: %v", err)
			}

			broadcastMsg := MessageItem{
				ID:               savedMsg.ID,
				ChatID:           savedMsg.ChatID,
				UserID:           savedMsg.UserID,
				Username:         client.Username,
				Text:             savedMsg.Text,
				IsEdited:         savedMsg.IsEdited,
				ReplyToMessageID: savedMsg.ReplyToMessageID,
				CreatedAt:        savedMsg.CreatedAt,
				Replies:          []MessageItem{},
			}

			broadcastData, _ := json.Marshal(Message{
				Type:    "message",
				Payload: broadcastMsg,
			})

			h.BroadcastToRoom <- &RoomMessage{
				RoomID: roomIDStr,
				Data:   broadcastData,
			}
		}
	}
}
