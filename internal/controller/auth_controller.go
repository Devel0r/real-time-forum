package controller

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

type AuthController struct {
	ARepo *repository.AuthRepository
}

func NewAuthController(db *sqlite.Database) *AuthController {
	return &AuthController{
		ARepo: repository.NewAuthRepository(db),
	}
}

// SignUpPage, GET
func (actl *AuthController) SignUpPage(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmpPath("signUp")))

	// TODO: status ok in the header response
	w.WriteHeader(http.StatusOK)

	// TODO: execute template
	if err := tmp.Execute(w, http.StatusOK); err != nil {
		slog.Error(err.Error())
		return
	}
}

// SignUp, POST
func (actl *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmpPath("signUp")))

    
	
}


// SignIn

// SignOut
