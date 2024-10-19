package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/Pruel/real-time-forum/internal/controller/wschat"
	"github.com/Pruel/real-time-forum/pkg/serror"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

// Chat struct
type ChatRepository struct {
	DB *sqlite.Database
}

// new_constructor
func NewChatReposotory(db *sqlite.Database) *ChatRepository {
	return &ChatRepository{
		DB: db,
	}
}

// SaveRoom
func (c *ChatRepository) SaveRoom(room *wschat.SRoom) (id string, err error) {
	if room == nil {
		return "", errors.New("error, nil struct pointer")
	}

	data, err := json.Marshal(room.Clients)
	if err != nil {
		return "", err
	}

	msgData, err := json.Marshal(room.LastMessage)
	if err != nil {
		return "", err
	}

	if err := c.DB.SQLite.QueryRow("SELECT id FROM rooms WHERE id=?", room.ID).Scan(&id); err != nil {
		if err != sql.ErrNoRows {
			return id, serror.ErrRoomAlreadyExists
		}
	}

	if _, err := c.DB.SQLite.Exec("INSERT INTO rooms(id, name, clients, client_creator_id, client_invited_id, last_msg) values(?, ?, ?, ?, ?, ?)", room.ID,
		room.Name, string(data), room.ClientCretorID, room.ClientInvitedID, string(msgData)); err != nil {
		return "", err
	}

	return room.ID, nil
}

// GetRoomByID
func (c *ChatRepository) GetRoomByID(id string) (sroom wschat.SRoom, err error) {
	if id == "" {
		return sroom, serror.ErrEmptyRoomID
	}

	clients := ""
	lastMsg := ""
	err = c.DB.SQLite.QueryRow("SELECT id, name, clients, client_creator_id, client_invited_id, last_msg FROM rooms WHERE id=?", id).Scan(&sroom.ID,
		&sroom.Name, &clients, &sroom.ClientCretorID, &sroom.ClientInvitedID, &sroom.LastMessage)
	if err != nil {
		return sroom, serror.ErrEmptyRoomID
	}

	if err := json.Unmarshal([]byte(clients), &sroom.Clients); err != nil {
		return sroom, err
	}

	if err := json.Unmarshal([]byte(lastMsg), &sroom.LastMessage); err != nil {
		return sroom, nil
	}

	return sroom, nil
}

// GetRooms
func (c *ChatRepository) GetAllRoomsByClientID(clientID string) (rooms []wschat.SRoom, err error) {
	if clientID == "" {
		return nil, errors.New("error, invalid room ids")
	}

	rrows, err := c.DB.SQLite.Query("SELECT id, name, clients, client_creator_id, client_invited_id, last_msg FROM rooms WHERE id=?")
	if err != nil {
		return nil, err
	}

	for rrows.Next() {
		room := wschat.SRoom{}
		cls := ""
		lastMsg := ""
		if err := rrows.Scan(&room.ID, &room.Name, &cls, &room.ClientCretorID, &room.ClientInvitedID, &lastMsg); err != nil {
			slog.Warn(err.Error())
			continue
		}

		if err := json.Unmarshal([]byte(cls), &room.Clients); err != nil {
			slog.Warn(err.Error())
			continue
		}

		if err := json.Unmarshal([]byte(lastMsg), &room.LastMessage); err != nil {
			slog.Warn(err.Error())
			continue
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

// DeleteRoomByID
func (c *ChatRepository) DeleteByRoomID(roomID string) error {
	if roomID == "" {
		return serror.ErrEmptyRoomID
	}

	if _, err := c.DB.SQLite.Exec("DELETE FROM rooms WHERE id=?", roomID); err != nil {
		return err
	}

	return nil
}

// SaveClients
func (c *ChatRepository) SaveClients(clients ...*wschat.SClient) (clients_ids []string, err error) {
	if clients == nil {
		return nil, errors.New("error, invalid clients data")
	}

	// ..Type -> func(1, 2, 3) -> []int{}
	// I have a questin

	for _, cl := range clients {
		id := ""
		if err := c.DB.SQLite.QueryRow("SELECT id FROM clients WHERE id=?", cl.ID).Scan(&id); err != nil {
			if err == sql.ErrNoRows { // TODO: ref: this code, DRY
				_, err := c.DB.SQLite.Exec("INSERT INTO clients(id, username, avatar, rooms_id) VALUES(?, ?, ?, ?)", cl.ID,
					cl.Username, cl.Avatar, cl.RoomID)
				if err != nil {
					return nil, err
				}
				clients_ids = append(clients_ids, cl.ID)
			}
		}

		if id != "" {
			slog.Warn(errors.New("warning, clients by this id already exists into dabatase").Error())
			clients_ids = append(clients_ids, id)
			continue
		}

		_, err := c.DB.SQLite.Exec("INSERT INTO clients(id, username, avatar, rooms_id) VALUES(?, ?, ?, ?)", cl.ID,
			cl.Username, cl.Avatar, cl.RoomID)
		if err != nil {
			return nil, err
		}
		clients_ids = append(clients_ids, cl.ID)
	}

	return clients_ids, nil
}

// GetClientByUsername
func (c *ChatRepository) GetClientByUsername(username string) (client wschat.SClient, err error) {
	if username == "" {
		return client, errors.New("error, clients username is empty")
	}

	err = c.DB.SQLite.QueryRow("SELECT id, username, avatar, rooms_id FROM clients WHERE username=?", username).Scan(&client.ID,
		&client.Username, &client.Avatar, &client.RoomID)
	if err != nil {
		return client, err
	}

	return client, nil
}

// GetAllChatUsers
func (c *ChatRepository) GetAllClients() (clients []wschat.SClient, err error) {
	crows, err := c.DB.SQLite.Query("SELECT id, username, avatar, rooms_id FROM clients")
	if err != nil {
		return nil, err
	}

	for crows.Next() {
		cl := &wschat.SClient{}
		if err := crows.Scan(cl.ID, cl.Username, cl.Avatar, cl.RoomID); err != nil {
			slog.Warn(err.Error())
			continue
		}
		clients = append(clients, *cl)
	}

	return clients, nil
}

// SaveMessage
func (c *ChatRepository) SaveMessage(msg *wschat.Message) (msgID string, err error) {
	if msg == nil {
		return "", errors.New("error, invalid message, msg is nil")
	}

	_, err = c.DB.SQLite.Exec("INSERT INTO messages(id, username, room_id, content, created_at, is_delivered, is_read) VALUES (?, ?, ?, ?, ?, ?, ?)",
		msg.ID, msg.From, msg.RoomID, msg.Content, msg.CreatedAt, msg.IsDelivered, msg.IsRead)
	if err != nil {
		return "", err
	}

	return msg.ID, nil
}

// GetAllMessagesByRoomID
func (c *ChatRepository) GetAllMessagesByRoomID(roomID string) (msgs []*wschat.Message, err error) {
	if roomID == "" {
		return nil, errors.New("error, roomID is empty")
	}

	mrow, err := c.DB.SQLite.Query("SELECT id, username, room_id, content, created_at, is_delivered, is_read FROM messages WHERE room_id=?", roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return msgs, err
		}
		return nil, err
	}

	for mrow.Next() {
		msg := wschat.Message{}
		if err := mrow.Scan(&msg); err != nil {
			slog.Warn(err.Error())
			continue
		}
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

// DeleteMessageByID
func (c *ChatRepository) DeleteMessageByID(msgID string) error {
	if msgID == "" {
		return errors.New("error, message id is empty")
	}

	if _, err := c.DB.SQLite.Exec("DELETE FROM messages WHERE id=?", msgID); err != nil {
		return err
	}

	return nil
}

// GetMoreMessagesByRoomID
