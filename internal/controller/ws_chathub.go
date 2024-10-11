package controller

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type ChatHub struct {
	IsGroupChat bool
	Rooms       map[string]*Room
	Clients     map[string]*Client
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan *Message
	Mu          *sync.RWMutex
	ChatRepo    *repository.ChatRepository // TODO: ChatRpepository
}

// chathub or room
type Room struct {
	ID      string             `sql:"id" json:"id"`
	Name    string             `sql:"name" json:"name"`
	Clients map[string]*Client `sql:"clients" json:"clients"`
}

type Client struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
	Conn     *websocket.Conn
	Message  chan *Message
}

type Message struct {
	ID          string    `sql:"id" json:"id"`
	From        string    `sql:"from" json:"from"`
	RoomID      string    `sql:"room_id" json:"room_id"`
	Content     string    `sql:"content" json:"content"`
	CreatedAt   time.Time `sql:"created_at" json:"created_at"`
	IsDelivered bool
	IsRead      bool
}

func NewChat(db *sqlite.Database) *ChatHub {
	return &ChatHub{
		Rooms:      make(map[string]*Room, 10),
		Clients:    make(map[string]*Client, 2),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 10),
		Mu:         &sync.RWMutex{},
		ChatRepo:   repository.NewChatReposotory(db),
	}
}

// Run
func (h *ChatHub) Run() {
	for {
		select {
		case cl := <-h.Register: // Client{ID: 12, RoomID: gophers, other:... }
			// 1. проверяем наличие такого чата
			if _, ok := h.Rooms[cl.RoomID]; ok {
				// 2. если такой чат есть, получаем его из мапы чатов для того чтобы добавить юзера
				chat := h.Rooms[cl.RoomID]

				// 3. добавляем юзера
				if _, ok := chat.Clients[cl.ID]; !ok {
					chat.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			// 1. check chat exists
			if _, ok := h.Rooms[cl.RoomID]; ok {
				// 2. if chat exists, check user exists in this chat room
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok { // [][]int{ []int{}}
					// 3. check count of users in this chat room, before unregister (delete) a user from this chat
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						// 4. broadcast the messaage for other users for information about the user left
						h.Broadcast <- &Message{
							ID:        GetUUID(),
							From:      cl.Username,
							RoomID:    cl.RoomID,
							Content:   fmt.Sprintf("User: %s left the %s\n", cl.Username, cl.RoomID),
							CreatedAt: time.Now(),
						}
					}
					// 5. delete this user from room | chat list users
					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					// 6. and after close the user message chan
					close(cl.Message)
				}
			}
		case msg := <-h.Broadcast:
			// 1. get rooms | chats by roomID from msg.RoomID
			if _, ok := h.Rooms[msg.RoomID]; ok {
				// 2. get every users | clients in this rooms
				for _, user := range h.Rooms[msg.RoomID].Clients {
					// 3. broadcast a message for every users
					user.Message <- msg
				}
			}
		}
	}
}

// getGoogleUUID
func GetUUID() (suuid string) {
	sguuid := uuid.DefaultGenerator
	uuidV4, err := sguuid.NewV4()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	suuid = uuidV4.String()

	return suuid
}

// WriteMessage
func (c *Client) WriteMessage() {
	// before close func, we need close web socket connection
	defer func() {
		c.Conn.Close()
	}()

	for {
		// 1. if we recieve a message from a user msg channel
		msg, ok := <-c.Message
		if !ok {
			slog.Error("Con't get a message from user message channel")
			return
		}
		// 3. send message -> web_socket
		c.Conn.WriteJSON(msg)
	}
}

// ReadMessage
func (c *Client) ReadMessage(ch *ChatHub) {
	// before close the func, we need close websocket conn, and unregister user from this room | chat
	defer func() {
		ch.Unregister <- c
		c.Conn.Close()
	}()

	for {
		// 1. if recive a new message from websocket con, we readd this msg
		_, msg, err := c.Conn.ReadMessage() // TODO: json
		// 2. if error occured while reading the message, log this err, and break infinity loop
		if err != nil { // errors.Is()
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error(err.Error())
			}
			break
		}

		// 3. create a new instance of thh message, and assign data
		message := &Message{
			ID:          GetUUID(),
			From:        c.Username,
			RoomID:      c.RoomID,
			Content:     string(msg),
			CreatedAt:   time.Now(),
			IsDelivered: true,
		}

		// 4. after broadcast this message to other users
		ch.Broadcast <- message
	}
}
