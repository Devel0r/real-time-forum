package controller

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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

func (actl *AuthController) ExecTmp(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	switch url {
	case "/sign-up":
		execTemplate(w, "signUp")
	case "/sign-in":
		execTemplate(w, "signIn")
	default:
		TmpPath = GetTmpPath("main")
	}
}

func execTemplate(w http.ResponseWriter, tmpPath string) error {
	tmp, err := template.ParseFiles(GetTmpPath(tmpPath))
	if err != nil {
		fmt.Println("Error, template: ", err)
	}

	w.WriteHeader(http.StatusOK)

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
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
	// get from native string password hash password
	hash, err := getPasswordHash(pass)
	if err != nil {
		slog.Warn(err.Error())
		actl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	user.PasswordHash = hash

	userID, err := actl.ARepo.SaveUser(&user)
	if err != nil {
		slog.Error(err.Error())
		actl.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slog.Debug("user saved into database with ", "id", userID)

	cookie, err := createCookie()
	if err != nil {
		slog.Error(err.Error())
		actl.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
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
		actl.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	// a user send request -> coockie
	// ++++++++++++++++++++++++++++++++++++++++++++++++++
	// coockie, err := r.Cookie("sessionID")
	// if err != http.ErrNoCookie {
	// 	slog.Error(err.Error())
	// 	a.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	// 	return
	// } // loading login page ---

	// sessionID := coockie.Value // user session_id == database session_id

	// TODO: get session_id from database, if we found the sesion with this session_id
	// TODO: if ok, a) create a new sesion coockie, and b) set this new session coockie to http response
	// TODO: redirect user to main page, with status code see other
	// TODO: compare session_id value beetwin coockie and session_id, user_id from sessions table (database)
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++

	// if coockie not exists
	loginEmail := r.FormValue("username_or_email")
	password := r.FormValue("password")

	if err := ValidateDateForLogin(loginEmail, password); err != nil {
		slog.Warn(err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusBadRequest)
		return
	}

	if !a.isValidUser(w, loginEmail, password) {
		slog.Warn("error, user try login with invalid password")
		http.Redirect(w, r, "/sign-in", http.StatusBadRequest)
		return
	}

	// create a new session
	coockie, err := createCookie()
	if err != nil {
		slog.Warn(err.Error())
		a.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.SetCookie(w, coockie)

	fmt.Println("User successful logined: ", coockie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (a *AuthController) isValidUser(w http.ResponseWriter, loginEmail string, password string) bool {
	if password == "" || loginEmail == "" {
		return false
	}

	var err error
	user := &model.User{}
	if sdata := strings.Split(loginEmail, "@"); len(sdata) == 2 {
		// email
		user, err = a.ARepo.GetUserByEmail(loginEmail, user)
		if err != serror.ErrUserNotFound {
			slog.Warn(err.Error())
			a.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return false
		}

	} else { // TODO: fix unused var compiler error
		user, err = a.ARepo.GetUserByUsername(loginEmail)
		if err != serror.ErrEmptyEmail {
			fmt.Println(err.Error())
			slog.Warn(err.Error()) 
			a.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return false
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return false
	}

	fmt.Println("User is valid: ", user)

	return true // TODO: fix error
}

func (a *AuthController) SignOut(w http.ResponseWriter, r *http.Request) {
	coockie, err := r.Cookie("sessionID")
	if err != nil {
		slog.Warn(err.Error())
		a.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// removeSessionByUUID
	// chcking error
	id, err := a.ARepo.RemoveSessionByUUID(coockie.Value)
	if err != nil {
		slog.Warn(err.Error())
		a.ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slog.Debug("successful remove session with", "id", id)

	coockie = &http.Cookie{}
	http.SetCookie(w, coockie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createCookie() (*http.Cookie, error) {
	expiresTime := time.Now().Add(time.Hour * 4)

	uuid := uuid.DefaultGenerator
	uuidV4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	value := uuidV4.String()
	// value = fmt.Sprintf("%d:%s", id, value) // id:session_id

	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    value,
		Expires:  expiresTime,
		HttpOnly: true,
		Secure:   false,
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

func ValidateDateForLogin(data, password string) error {
	if data == "" || password == "" {
		return serror.ErrEmptyFieldLogin
	}

	if sdata := strings.Split(data, "@"); len(sdata) == 2 {
		if ok := validator.ValidateEmail(data); !ok {
			return serror.ErrInvalidEmail
		}
	}

	// prefixs <script>, sql script - sql injection

	return nil
}

// Blume Filter
