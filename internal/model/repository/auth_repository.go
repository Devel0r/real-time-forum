package repository

import (
	"database/sql"
	"errors"

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
// check session,

// SignOut
// remove session,

// GetUserByUsername
func (a *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	if username == "" {
		return nil, serror.ErrEmptyUsername
	}

	user := model.User{}
	if err := a.DB.SQLite.QueryRow("SELECT id, login, age, gender, name, surname, email, password_hash FROM users WHERE login=?", username).Scan(&user.Id,
		&user.Login, &user.Age, &user.Gender, &user.Name, &user.Surname, &user.Email, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, serror.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail
func (a *AuthRepository) GetUserByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, serror.ErrEmptyEmail
	}

	user := model.User{}
	if err := a.DB.SQLite.QueryRow("SELECT id, login, age, gender, name, surname, email, password_hash FROM users WHERE email=?", email).Scan(&user.Id,
		&user.Login, &user.Age, &user.Gender, &user.Name, &user.Surname, &user.Email, &user.PasswordHash); err != nil {
		if err != sql.ErrNoRows {
			return nil, serror.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// SaveUser
func (a *AuthRepository) SaveUser(user *model.User) (id int, err error) {
	if user == nil {
		return 0, serror.ErrEmptyUserData
	}

	res, err := a.DB.SQLite.Exec("INSERT INTO users(login, age, gender, name, surname, email, password_hash) VALUES(?, ?, ?, ?, ?, ?, ?)",
		user.Login, user.Age, user.Gender, user.Name, user.Surname, user.Email, user.PasswordHash)
	if err != nil {
		return 0, err
	}

	ID, err := res.LastInsertId()
	return int(ID), err
}

func (a *AuthRepository) SaveCookie(session *model.Session) (id int, err error) {
	if session == nil {
		return 0, serror.ErrEmptyCookieData
	}

	res, err := a.DB.SQLite.Exec("INSERT INTO sessions(id, user_id, expires_at, created_at) VALUES(?, ?, ?, ?)",
		session.Id, session.UserId, session.ExpiredAt, session.CreatedAt)
	if err != nil {
		return 0, err
	}

	ID, err := res.LastInsertId()
	return int(ID), err
}

// Save cookie по факту тоже самое что сверху но только значение передаём другие, передаём поля sql.session

// removeSessionByUUID
func (a *AuthRepository) RemoveSessionByUUID(uuid string) (int, error) {
	if uuid == "" {
		return 0, errors.New("error, incorrenct session uuid")

	}

	res, err := a.DB.SQLite.Exec("DELETE FROM sessions WHERE id=?", uuid)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()

	return int(id), nil
}
