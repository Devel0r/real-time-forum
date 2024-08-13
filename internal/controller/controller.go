package controller

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

// Main Controller
type Controller struct {
	*AuthController
	// category controller
	// post controller
	// comment controller
}

// New return a new instance of the Controller
func New(db *sqlite.Database) *Controller {
	return &Controller{
		AuthController: NewAuthController(db),
	}
}

// GetWd
func GetWd() (wd string) {
	wd, _ = os.Getwd()
	return wd
}

// GetTmpPath
func GetTmpPath(tmpName string) (tmpPath string) {
	switch tmpName {
	case "signUp", "sign_up", "signUp.html":
		tmpPath = GetWd() + "internal/view/template/sign_up.html"
	case "signIn", "sign_in", "signIn.html":
		tmpPath = GetWd() + "internal/view/template/sign_in.html"
	}

	return tmpPath
}

// Controller of the main page
func (ctl *Controller) MainController(w http.ResponseWriter, r *http.Request) {
	slog.Debug("main page")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Main page"))
}
