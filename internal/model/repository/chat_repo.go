package repository

import "github.com/Pruel/real-time-forum/pkg/sqlite"

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

// methods
func (c *ChatRepository) AmongsAss() {

}

// SavePvChat

// DeletePvChat

// GetAllChatUsers

// SaveMessage

// DeleteMessage

// GetMoreMessagesByRoomID
