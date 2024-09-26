package controller

import (
	"fmt"
	"html/template"

	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

type Controller struct {
	*AuthController
}

// New return a new instance of the Controller
func New(db *sqlite.Database) *Controller {
	return &Controller{
		AuthController: NewAuthController(db),
	}
}

func (ctl *Controller) GetStaticPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return "./view/static"
	}
	return filepath.Join(wd, "view", "static")
}

func (ctl *Controller) MainController(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/static") {
		http.NotFound(w, r)
		return
	}
	// Путь к index.html
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	indexPath := filepath.Join(wd, "view", "template", "index.htm")
	fmt.Printf("\n\n Path to index.html: %s \n\n", indexPath)

	file, err := os.Open(indexPath)
	defer file.Close()
	if err != nil {
		fmt.Printf("\nnError: %s\n\n", err)
	}
	fmt.Println("file name: ", file.Name())

	tmp := template.Must(template.ParseFiles(indexPath))

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
	}
}
