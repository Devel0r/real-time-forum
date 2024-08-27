package controller

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Pruel/real-time-forum/internal/model"
	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/serror"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
	"github.com/Pruel/real-time-forum/pkg/validator"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	ARepo *repository.AuthRepository
}

func NewAuthController(db *sqlite.Database) *AuthController {
	return &AuthController{
		ARepo: repository.NewAuthRepository(db),
	}
}

// SignUpPage
func (actl *AuthController) SignUpPage(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmpPath("signUp")))

	w.WriteHeader(http.StatusOK)

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
		return
	}
}

// SignUp
func (actl *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {

	if err := actl.validateUserData(r); err != nil {
		slog.Warn(err.Error())
		actl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
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
	
	hash, err := getPasswordHash(pass)
	if err != nil {
		slog.Warn(err.Error())
		actl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	user.PasswordHash = hash

	id, err := actl.ARepo.SaveUser(&user)
	if err != nil {
		slog.Error(err.Error())
		actl.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slog.Debug("user saved into database with ", "id", id)

	cookie, err := createCookie(id)
	if err != nil {
		slog.Error(err.Error())
		actl.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	http.SetCookie(w, cookie)

	fmt.Println("Cookie: ", cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *AuthController) SignIn(w http.Response, r *http.Request) {
	// Пользователь должен иметь возможность подключаться, используя либо ник, либо e-mail в сочетании с паролем.
	
}

func (a *AuthController) SignOut(w http.Response, r http.Request) {

}

/*
	POST https://game.org https/1.1
	host: game.org
	other headers
	cookie: session_id:1eqweqweqweasdadsf + https + TLS
	login: daniil
	password: qwerty
	
 */
func createCookie(id int) (*http.Cookie, error) {
	expiresTime := time.Now().Add(time.Hour * 4) // time.Duration, time.Time = 22:15:16 28.08.2024

	uuid := uuid.DefaultGenerator
	uuidV4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	value := uuidV4.String()
	value = fmt.Sprintf("%d:%s", id, value)

	cookie := &http.Cookie {
		Name: "sessionID",
		Value: value, 
		Expires: expiresTime,  
		HttpOnly: true,
		Secure: false,
	} 

	return cookie, nil
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
