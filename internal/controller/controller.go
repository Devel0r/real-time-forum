package controller

import (
	"fmt"
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
		return "./internal/view/static"
	}
	return filepath.Join(wd, "internal", "view", "static")
}

func (ctl *Controller) MainController(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/static/") {
		http.NotFound(w, r)
		slog.Warn("Page not found: " + r.URL.Path)
		return
	}

	// Путь к index.html
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.Error("Failed to get working directory: " + err.Error())
		return
	}
	indexPath := filepath.Join(wd, "internal", "view", "template", "index.html")
	fmt.Printf("\n\n Path to index.html: %s \n\n", indexPath)

	http.ServeFile(w, r, indexPath)
	slog.Info("Successfully served index.html")
}
