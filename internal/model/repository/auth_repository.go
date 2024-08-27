package repository

import (
	"database/sql"
	"fmt"

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
	if err := a.DB.SQLite.QueryRow("SELECT (login) FROM users WHERE login=?", username).Scan(user.Login); err != nil {
		if err == sql.ErrNoRows {
			return nil, serror.ErrUserNotFound
		}
		fmt.Println("unknown error")
		return nil, err
	}

	return user, nil
}

// SaveUser
func (a *AuthRepository) SaveUser(user *model.User) (id int, err error) {
	if user == nil {
		return 0, err
	}

	res, err := a.DB.SQLite.Exec("INSERT INTO users(login, age, gender, name, surname, email, password_hash) VALUES(?, ?, ?, ?, ?, ?, ?)",
		user.Login, user.Age, user.Gender, user.Name, user.Surname, user.Email, user.PasswordHash)
	if err != nil {
		return 0, err
	}

	ID, err := res.LastInsertId()
	return int(ID), err
}


// Save cookie по факту тоже самое что сверху но только значение передаём другие, передаём поля sql.session 
