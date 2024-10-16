package repository

import (
	"encoding/json"
	"errors"

	"github.com/Pruel/real-time-forum/internal/controller/wschat"
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

	ids := []string{}

	for _, cl := range room.Clients {
		ids = append(ids, cl.ID)
	}

	data, err := json.Marshal(ids)
	if err != nil {
		return "", err
	}

	_, err = c.DB.SQLite.Exec("INSERT INTO rooms(id, name, clients_id) values(?, ?, ?)", room.ID,
		room.Name, string(data))
	if err != nil {
		return "", err
	}

	return room.ID, nil
}

// GetRoomByID

// GetRooms

// DeleteRoomByID

// GetAllChatUsers

// SaveMessage

// DeleteMessage

// GetMoreMessagesByRoomID
