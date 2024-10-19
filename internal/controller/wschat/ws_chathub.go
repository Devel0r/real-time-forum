package wschat

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type ChatHub struct {
	IsGroupChat bool
	Rooms       map[string]*Room
	Clients     map[string]*Client
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan Message
	Mu          *sync.RWMutex
}

// chathub or room
type Room struct {
	ID      string             `sql:"id" json:"room_id"`
	Name    string             `sql:"name" json:"name"`
	Clients map[string]*Client `sql:"clients" json:"clients"`
}

type Client struct {
	ID       string `json:"client_id"`
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
	Conn     *websocket.Conn
	Message  chan Message
}

type Message struct {
	ID          string    `sql:"id" json:"message_id"`
	From        string    `sql:"from" json:"from"`
	RoomID      string    `sql:"room_id" json:"room_id"`
	Content     string    `sql:"content" json:"content"`
	CreatedAt   time.Time `sql:"created_at" json:"created_at"`
	IsDelivered bool
	IsRead      bool
}

// SimpleClient, SimpleRoom client struct clone without unsupported fields (json)
type SClient struct {
	ID       string `json:"client_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	RoomID   string `json:"room_id"`
	IsOnline bool   `json:"is_online"`
}

type SRoom struct {
	ID              string             `json:"room_id"`
	Name            string             `json:"name"`
	Clients         map[string]SClient `json:"clients"`
	ClientCretorID  string             `json:"client_cretor_id"`
	ClientInvitedID string             `json:"client_invited_id"`
	LastMessage     *Message           `json:"last_message"`
}

func NewChat() *ChatHub {
	return &ChatHub{
		Rooms:      make(map[string]*Room, 10),
		Clients:    make(map[string]*Client, 2),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message, 10),
		Mu:         &sync.RWMutex{},
	}
}

// Run
func (h *ChatHub) Run() {
	for {
		select {
		case cl := <-h.Register: // Client{ID: 12, RoomID: gophers, other:... }
			// 1. check for exists this chat | room
			slog.Debug(fmt.Sprintf("request for register the %s for the %s chat\n", cl.Username, cl.RoomID))
			if _, ok := h.Rooms[cl.RoomID]; ok {
				// 2. if this chat | room exists, select this chat for regisring this client on this chat
				chat := h.Rooms[cl.RoomID]
				slog.Debug("register channel, chat exists") // TODO: remove this log

				// 3. add this client
				if _, ok := chat.Clients[cl.ID]; !ok {
					chat.Clients[cl.ID] = cl
					slog.Info(fmt.Sprintf("the %s successful registered in the %s chat\n", cl.Username, cl.RoomID))
				}
			}
		case cl := <-h.Unregister:
			// 1. check chat exists
			slog.Debug("request for unregister the %s from the %s\n", cl.Username, cl.RoomID)
			if _, ok := h.Rooms[cl.RoomID]; ok {
				// 2. if chat exists, check user exists in this chat room
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok { // [][]int{ []int{}}
					// 3. check count of users in this chat room, before unregister (delete) a user from this chat
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						// 4. broadcast the messaage for other users for information about the user left
						h.Broadcast <- Message{
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
					slog.Info(fmt.Sprintf("the %s successful unregistered from the %s\n", cl.Username, cl.RoomID))
				}
			}
		case msg := <-h.Broadcast:
			// 1. get rooms | chats by roomID from msg.RoomID
			slog.Debug("recieve a message in broadcast channel")
			// slog.Debug(fmt.Sprintf("broadcast:\n\tmsg: %s\n\tfrom:%s\n\for the chat:%s\n", msg.Content, msg.From, msg.RoomID))
			if _, ok := h.Rooms[msg.RoomID]; ok {
				// 2. get every users | clients in this rooms
				for _, cl := range h.Rooms[msg.RoomID].Clients {
					// 3. broadcast a message for every users
					slog.Debug("message successful broadcast for every this chat clients", "client", cl.Username, "msgID", msg.ID)
					cl.Message <- msg
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
		slog.Info("client.WriteMessage, successful written the msg", "msg", msg)
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
		_, msg, err := c.Conn.ReadMessage()
		// 2. if error occured while reading the message, log this err, and break infinity loop
		if err != nil { // errors.Is()
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error(err.Error())
				break
			}
			break
		}

		// 3. create a new instance of thh message, and assign data
		message := Message{
			ID:          GetUUID(),
			From:        c.Username,
			RoomID:      c.RoomID,
			Content:     string(msg),
			CreatedAt:   time.Now(),
			IsDelivered: true,
		}

		// 4. after broadcast this message to other users
		fmt.Println("client.ReadMessage, send msg to broadcast channel") // TODO: remove this fmt log
		ch.Broadcast <- message
	}

	slog.Debug("client.ReadMessage, the ws connection was closed")
}
