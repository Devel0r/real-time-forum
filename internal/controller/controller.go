package controller

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

var TmpPath string

type TemplateData struct {
	Username string
}

// Main Controller
type Controller struct {
	*AuthController
	*PostController
	*CommentController
	// category controller for Admin set up maybe in future
}

func New(db *sqlite.Database) *Controller {
	return &Controller{
		AuthController: NewAuthController(db),
	}
}

func (ctl *Controller) MainController(w http.ResponseWriter, r *http.Request) {

	tmp := template.Must(template.ParseFiles(GetTmpPath("index")))

	// cookie = session - sessionUUID
	userID, err := ctl.AuthController.ARepo.GetUserIDFromSession(w, r)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	name, err := ctl.AuthController.ARepo.GetUserNameByUserID(userID)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	data := TemplateData{
		Username: name,
	}

	if err := tmp.Execute(w, data); err != nil {
		slog.Error(err.Error())
		return
	}
}

func GetWd() (wd string) {
	wd, _ = os.Getwd()
	return wd
}

func GetTmpPath(tmpName string) (tmpPath string) {
	switch tmpName {
	case "signUp":
		tmpPath = GetWd() + "/internal/view/template/sign_up.html"
	case "signIn":
		tmpPath = GetWd() + "/internal/view/template/sign_in.html"
	case "post":
		tmpPath = GetWd() + "/internal/view/template/post.html"
	case "error":
		tmpPath = GetWd() + "/internal/view/template/error.html"
	case "index":
		tmpPath = GetWd() + "/internal/view/template/index.html"
	}

	return tmpPath
}

func (actl *AuthController) ExecTmp(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	switch url {
	case "/sign-up":
		execTemplate(w, "signUp")
	case "/sign-in":
		execTemplate(w, "signIn")
	case "/post-create":
		execTemplate(w, "post")
	case "/":
		execTemplate(w, "index")
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
