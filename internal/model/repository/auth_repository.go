package repository

import (
	"database/sql"

	"github.com/Pruel/real-time-forum/internal/model"
	"github.com/Pruel/real-time-forum/pkg/serror"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

type AuthRepository struct {
	DB *sqlite.Database
}

func NewAuthRepository(db *sqlite.Database) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

// SignUp
// create session

// SignIn
// check session

// SignOut
// remove session


// GetUserByUsername
func (a *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	if username == "" {
		return nil, serror.ErrEmptyUsername
	}

	user := &model.User{}
	if err := a.DB.SQLite.QueryRow("SELECT * FROM users WHERE username=?", username).Scan(user); err != nil {
		if err == sql.ErrNoRows {
			return nil, serror.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}