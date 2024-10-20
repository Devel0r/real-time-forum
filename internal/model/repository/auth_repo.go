package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

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
func (a *AuthRepository) GetUserIdByUsername(username string) (int, error) {
	if username == "" {
		return 0, errors.New("username cannot be empty")
	}

	var userId int
	query := `SELECT id FROM users WHERE login = ?`

	// Выполняем SQL-запрос и сканируем данные в userId
	err := a.DB.SQLite.QueryRow(query, username).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, serror.ErrUserNotFound // Если пользователь не найден
		}
		return 0, err // Возвращаем другую ошибку
	}

	return userId, nil // Возвращаем ID пользователя
}

// GetUserByUsername
func (a *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	if username == "" {
		return nil, serror.ErrEmptyUsername
	}

	user := model.User{}
	err := a.DB.SQLite.QueryRow("SELECT id, login, age, gender, name, surname, email, password_hash FROM users WHERE login=?", username).Scan(
		&user.Id,
		&user.Login,
		&user.Age,
		&user.Gender,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serror.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail
func (a *AuthRepository) GetUserByOnlyEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, serror.ErrEmptyEmail
	}

	user := model.User{}
	// Выполняем SQL-запрос и сканируем данные в поля структуры user
	err := a.DB.SQLite.QueryRow("SELECT id, login, age, gender, name, surname, email, password_hash FROM users WHERE email=?", email).Scan(
		&user.Id, &user.Login, &user.Age, &user.Gender, &user.Name, &user.Surname, &user.Email, &user.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователь не найден, возвращаем ErrUserNotFound
			return nil, serror.ErrUserNotFound
		}
		// Возвращаем другие возможные ошибки
		return nil, err
	}

	// Возвращаем найденного пользователя
	return &user, nil
}

func (a *AuthRepository) GetUserByEmail(email string, user *model.User) (*model.User, error) {
	if email == "" {
		return nil, serror.ErrEmptyEmail
	}

	// Выполняем SQL-запрос и сканируем данные в поля структуры user
	err := a.DB.SQLite.QueryRow("SELECT login, age, gender, name, surname, email, password_hash FROM users WHERE email=?", email).Scan(
		&user.Login, &user.Age, &user.Gender, &user.Name, &user.Surname, &user.Email, &user.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователь не найден, возвращаем ErrUserNotFound
			return nil, serror.ErrUserNotFound
		}
		// Возвращаем другие возможные ошибки
		return nil, err
	}

	// Возвращаем найденного пользователя
	return user, nil
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

// removeSessionByUUID
func (a *AuthRepository) RemoveSessionByUUID(uuid string) (id int, err error) {
	if uuid == "" {
		return 0, errors.New("error, incorrenct session uuid")

	}

	res, err := a.DB.SQLite.Exec("DELETE FROM sessions WHERE id=?", uuid)
	if err != nil {
		return 0, err
	}

	ID, _ := res.LastInsertId()

	return int(ID), nil
}

func (a *AuthRepository) GetUserIDFromSession(w http.ResponseWriter, r *http.Request) (int, error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		if err == http.ErrNoCookie {
			slog.Warn("Unauthorized user", "error", err.Error())
			http.Redirect(w, r, "/sign-in", http.StatusForbidden)
			return 0, err
		}
	}

	sessionUUID := cookie.Value
	user := model.User{}
	if err = a.DB.SQLite.QueryRow("SELECT user_id FROM sessions WHERE id=?", sessionUUID).Scan(&user.Id); err != nil {
		return 0, err
	}

	return user.Id, nil
}

// GetUserByUserID
func (a *AuthRepository) GetUserByUserID(userID int) (*model.User, error) {
	if userID == 0 {
		return nil, serror.ErrEmptyUserData
	}

	user := model.User{}
	err := a.DB.SQLite.QueryRow("SELECT id, login, age, gender, name, surname, email, password_hash FROM users WHERE id=?", userID).Scan(
		&user.Id, &user.Login, &user.Age, &user.Gender, &user.Name, &user.Surname, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAllUsers
func (a *AuthRepository) GetAllUsers() ([]model.User, error) {
	urows, err := a.DB.SQLite.Query("SELECT id, login, age, gender, name, surname, email, password_hash, rooms_id FROM users")
	if err != nil {
		return nil, err  
	}
	
	users := []model.User{}
	for urows.Next() {
		user :=  model.User{}
		rooms :=  ""
		err := urows.Scan(&user.Id, &user.Login, &user.Age, &user.Gender, &user.Name, &user.Surname, &user.Email, &user.PasswordHash, &rooms)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(rooms), &user.RoomsID); err != nil {
			return nil,  err
		}
		
		users = append(users, user)
	}
	
	return users, nil
}
