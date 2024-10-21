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

	msgData, err := json.Marshal(room.LastMessage)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	if err := c.DB.SQLite.QueryRow("SELECT id FROM rooms WHERE id=?", room.ID).Scan(&id); err != nil {
		if err != sql.ErrNoRows {
			return id, serror.ErrRoomAlreadyExists
		}
	}

	if _, err := c.DB.SQLite.Exec("INSERT INTO rooms(id, name, client_creator_id, client_invited_id, last_msg) values(?, ?, ?, ?, ?)", room.ID,
		room.Name, room.ClientCretorID, room.ClientInvitedID, string(msgData)); err != nil {
		return "", err
	}

	// update the clients rooms list
	// get the old rooms list
	usersIDs := []int{room.ClientCretorID, room.ClientInvitedID}
	for _, userID := range usersIDs {
		rms := ""
		if err := c.DB.SQLite.QueryRow("SELECT rooms_id FROM users WHERE id=?", userID).Scan(&rms); err != nil {
			if err != sql.ErrNoRows {
				return "", err
			}
		}

		// add last created clinet room id in the old rooms list of the client
		jrooms := []string{}
		if err := json.Unmarshal([]byte(rms), &jrooms); err != nil {
			slog.Error(err.Error())
			return "", err
		}
		jrooms = append(jrooms, room.ID)

		nrData, err := json.Marshal(jrooms)
		if err != nil {
			slog.Error(err.Error())
			return "", err
		}

		// update the rooms list of the client
		if _, err := c.DB.SQLite.Exec("UPDATE users SET rooms_id=? WHERE id=?", string(nrData), userID); err != nil {
			return "", err
		}
	}

	return room.ID, nil
}

// AddNewUserRoomID

// GetRoomByID
func (c *ChatRepository) GetRoomByID(id string) (sroom wschat.SRoom, err error) {
	if id == "" {
		return sroom, serror.ErrEmptyRoomID
	}

	lastMsg := ""
	err = c.DB.SQLite.QueryRow("SELECT id, name, client_creator_id, client_invited_id, last_msg FROM rooms WHERE id=?", id).Scan(&sroom.ID,
		&sroom.Name, &sroom.ClientCretorID, &sroom.ClientInvitedID, &lastMsg)
	if err != nil {
		return sroom, err
	}

	if err := json.Unmarshal([]byte(lastMsg), &sroom.LastMessage); err != nil {
		return sroom, nil
	}

	return sroom, nil
}

// SaveLastRoomMsg
func (c *ChatRepository) SaveLasgRoomMsg(roomID string, msg *wschat.Message) error {
	if roomID == "" || msg == nil {
		return errors.New("error, invalid args or nil struct pointer")
	}

	mdata, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := c.DB.SQLite.Exec("UPDATE rooms SET last_msg=? WHERE id=?", string(mdata), roomID); err != nil {
		return err
	}

	return nil
}

// GetRooms
func (c *ChatRepository) GetAllRoomsByClientID(clientID int) (rooms []wschat.SRoom, err error) {

	rrows, err := c.DB.SQLite.Query("SELECT id, name, client_creator_id, client_invited_id, last_msg FROM rooms WHERE client_creator_id=? OR client_invited_id=?", clientID, clientID)
	if err != nil {
		return nil, err
	}

	for rrows.Next() {
		room := wschat.SRoom{}
		lastMsg := ""
		if err := rrows.Scan(&room.ID, &room.Name, &room.ClientCretorID, &room.ClientInvitedID, &lastMsg); err != nil {
			if err == sql.ErrNoRows {
				slog.Warn(err.Error())
				continue
			}
			return nil, err
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
func (c *ChatRepository) GetAllMessagesByRoomID(roomID string) (msgs []wschat.Message, err error) {
	if roomID == "" {
		return nil, errors.New("error, roomID is empty")
	}

	mrow, err := c.DB.SQLite.Query("SELECT id, username, room_id, content, created_at, is_delivered, is_read FROM messages WHERE room_id=?", roomID)
	if err != nil {
		if err != sql.ErrNoRows {
			return msgs, err
		}
	}

	for mrow.Next() {
		msg := wschat.Message{}
		if err := mrow.Scan(&msg.ID, &msg.From, &msg.RoomID, &msg.Content, &msg.CreatedAt, &msg.IsDelivered, &msg.IsRead); err != nil {
			slog.Warn(err.Error())
			continue
		}
		msgs = append(msgs, msg)
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
