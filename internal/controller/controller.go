package controller

import (
	"fmt"
	// "html/template"
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
		slog.Warn("Page not found")
		return
	}
	// Путь к index.html
	indexPath := filepath.Join("internal", "view", "template", "index.html")
	fmt.Printf("\n\n Path to index.html: %s \n\n", indexPath)

	http.ServeFile(w, r, indexPath)
	slog.Info("Successful serve the index page file")

	// tmp := template.Must(template.ParseFiles(indexPath))
	// if err := tmp.Execute(w, nil); err != nil {
	// 	slog.Error(err.Error())
	// }
}
