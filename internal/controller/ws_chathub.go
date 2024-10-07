package controller

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ChatHub struct {
	ID          int
	Name        string
	IsGroupChat bool
	Clients     map[string]*Client
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan *Message
	Mu          *sync.RWMutex
	// ChatRepo	*repository.ChatRepository // TODO: ChatRpepository
}

// chathub or room
type Room struct {
	ID string	`sql:"id" json:"id"`
	Name string	`sql:"name" json:"name"`
	Clients map[string]*Client	`sql:"clients" json:"clients"`
}

type Client struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
	Conn     *websocket.Conn
	Message  chan *Message
}

type Message struct {
	ID          int       `sql:"id" json:"id"`
	From        string    `sql:"from" json:"from"`
	RoomID      string    `sql:"room_id" json:"room_id"`
	Content     string    `sql:"content" json:"content"`
	CreatedAt   time.Time `sql:"created_at" json:"created_at"`
	IsDelivered bool
	IsRead      bool
}
