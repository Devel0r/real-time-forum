package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Pruel/real-time-forum/internal/model"
	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/serror"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
	"github.com/Pruel/real-time-forum/pkg/validator"
)

type AuthController struct {
	ARepo *repository.AuthRepository
}

func NewAuthController(db *sqlite.Database) *AuthController {
	return &AuthController{
		ARepo: repository.NewAuthRepository(db),
	}
}

// SignUp
func (actl *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {

	if err := actl.validateUserData(r); err != nil {
		slog.Warn(err.Error())
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// получаем данные юзера
	user := model.User{}
	user.Login = r.FormValue("login")
	sAge := r.FormValue("age")
	user.Age, _ = strconv.Atoi(sAge)
	user.Gender = r.FormValue("gender")
	user.Name = r.FormValue("first_name")
	user.Surname = r.FormValue("last_name")
	user.Email = r.FormValue("email")

	pass := r.FormValue("password")
	// get from native string password hash password
	hash, err := getPasswordHash(pass)
	if err != nil {
		slog.Warn(err.Error())
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	user.PasswordHash = hash

	userID, err := actl.ARepo.SaveUser(&user)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slog.Debug("user saved into database with ", "id", userID)

	cookie, err := createCookie()
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// SaveCookie
	session := &model.Session{
		Id:        cookie.Value,
		UserId:    userID,
		ExpiredAt: cookie.Expires,
		CreatedAt: time.Now(),
	}

	_, err = actl.ARepo.SaveCookie(session)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	loginOrEmail := r.FormValue("username_or_email")
	password := r.FormValue("password")

	if err := ValidateDateForLogin(loginOrEmail, password); err != nil {
		slog.Warn(err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusBadRequest)
		return
	}

	if !a.isValidUser(w, loginOrEmail, password) {
		slog.Warn("error, user try login with invalid password\n")
		http.Redirect(w, r, "/sign-in", http.StatusBadRequest)
		return
	}

	userID, err := a.ARepo.GetUserIdByUsername(loginOrEmail)
	if err != nil {
		slog.Warn(err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusBadRequest)
		return
	}

	cookie, err := createCookie()
	if err != nil {
		slog.Warn(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	session := &model.Session{
		Id:        cookie.Value,
		UserId:    userID,
		ExpiredAt: cookie.Expires,
		CreatedAt: time.Now(),
	}

	_, err = a.ARepo.SaveCookie(session)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	http.SetCookie(w, cookie)

	slog.Info("User successfully logged")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (a *AuthController) SignOut(w http.ResponseWriter, r *http.Request) {
	coockie, err := r.Cookie("sessionID")
	if err != nil {
		slog.Warn(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	id, err := a.ARepo.RemoveSessionByUUID(coockie.Value)
	if err != nil {
		slog.Warn(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slog.Debug("successful remove session with", "id", id)

	coockie = &http.Cookie{}
	http.SetCookie(w, coockie)

	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

func createCookie() (*http.Cookie, error) {
	expiresTime := time.Now().Add(time.Hour * 4)

	uuid := uuid.DefaultGenerator
	uuidV4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	value := uuidV4.String()

	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    value,
		Expires:  expiresTime,
		HttpOnly: true,
		Secure:   false,
	}

	return cookie, nil
}

func (a *AuthController) isValidUser(w http.ResponseWriter, loginOrEmail string, password string) bool {
	if password == "" || loginOrEmail == "" {
		return false
	}

	var err error
	user := &model.User{}
	if sdata := strings.Split(loginOrEmail, "@"); len(sdata) == 2 {
		// email
		user, err = a.ARepo.GetUserByEmail(loginOrEmail, user)

		if err != nil {
			if err == serror.ErrUserNotFound {
				// Пользователь не найден
				slog.Warn("User not found")
				return false
			}
			// Другие ошибки
			slog.Warn(err.Error())
			ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return false
		}
	} else {
		user, err = a.ARepo.GetUserByUsername(loginOrEmail)
		if err != nil {
			if err == serror.ErrUserNotFound {
				// Пользователь не найден
				slog.Warn("User not found")
				return false
			}
			// Другие ошибки
			slog.Warn(err.Error())
			ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return false
		}
	}

	if user == nil {
		// Дополнительная проверка на nil
		slog.Warn("User is nil")
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		slog.Warn("Invalid password")
		return false
	}

	fmt.Println("User is valid: ", user)
	return true
}

func getPasswordHash(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (actl *AuthController) validateUserData(r *http.Request) error {
	minAge := 3
	maxAge := 110

	username := r.FormValue("login")
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	intAge, err := strconv.Atoi(age)
	if err != nil {
		return err
	}

	if intAge < minAge && intAge > maxAge {
		return serror.ErrIncorrectAge
	}

	if firstName == "" || lastName == "" || gender == "" {
		return serror.ErrIncorrectNameOrGender
	}

	if ok := validator.ValidateEmail(email); !ok {
		return serror.ErrInvalidEmail
	}

	if ok := validator.ValidatePassword(password); !ok {
		return serror.ErrInvalidPassword
	}

	if _, err := actl.ARepo.GetUserByUsername(username); err != nil {
		if err != serror.ErrUserNotFound {
			return err
		}
	}

	return nil
}

func ValidateDateForLogin(data, password string) error {
	if data == "" || password == "" {
		return serror.ErrEmptyFieldLogin
	}

	if sdata := strings.Split(data, "@"); len(sdata) == 2 {
		if ok := validator.ValidateEmail(data); !ok {
			return serror.ErrInvalidEmail
		}
	}

	return nil
}
