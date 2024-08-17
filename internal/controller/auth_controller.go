package controller

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

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
	tmp := template.Must(template.ParseFiles(GetTmpPath("signUp")))

	// TODO: status ok in the header response
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)

		// TODO: execute template
		if err := tmp.Execute(w, nil); err != nil {
			slog.Error(err.Error())
			return
		}
	}

	// TODO: prepate user instance, assign data

	// TODO: hashing password, bcrypt, uuid

	// TODO: save into database

	// TODO: create a new cookie for saving user session data

	// TODO: set this session cookie

	// TODO: redirect to main page
}

// SignIn

// SignOut

// validateUserData return nil if all data is valid, else return special error
func (actl *AuthController) validateUserData(r *http.Request) error {
	minAge := 3
	maxAge := 110

	// TODO: recieve data from frontend
	username := r.FormValue("username")    // GET (url), POST post form
	age := r.FormValue("age")              // - diapazon, 3 - 110
	gender := r.FormValue("gender")        // only two gender,
	firstName := r.FormValue("first_name") // "",
	lastName := r.FormValue("sur_name")    // ""
	email := r.FormValue("email")          // regexp, pattern, word (word or numbers + "@" + word + "." + "word")
	password := r.FormValue("password")    // a, A, 8, *, and 8 simbols

	// TODO: validate this date
	// request to database -> users -> user with this userename -> error -> message: try again with other username
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

	// TODO: view this code
	if _, err := actl.ARepo.GetUserByUsername(username); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	return nil
}
