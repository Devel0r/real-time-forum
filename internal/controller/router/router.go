package router

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Pruel/real-time-forum/internal/controller"
)

type Router struct { // public, protected, private -> Router
	Mux *http.ServeMux
	Ctl *controller.Controller
	// Chi *chi.Router
	// gorrilaMux *gorilla.Router
}

// New create and return a new instance of the ServeMux router
func New(ctl *controller.Controller) *Router {
	return &Router{
		Mux: http.NewServeMux(),
		Ctl: ctl,
	}
}

func (r *Router) InitRouter() {

	// TODO: inpmlement auth middleware for check every request -

	// statics

	wd, err := os.Getwd()
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	// r.Mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(wd+"/internal/view/static/css/"))))
	r.Mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(wd+"/internal/view/static/")))) // wd /

	r.Mux.HandleFunc("GET /", r.Ctl.MainController) // template -> router -> controller -> model  -> repository -> database

	// auth routes, sign-up, sign-in, sign-out
	r.Mux.HandleFunc("GET /sign-up", r.Ctl.AuthController.ExecTmp)
	r.Mux.HandleFunc("POST /sign-up", r.Ctl.AuthController.SignUp)
	r.Mux.HandleFunc("GET /sign-in", r.Ctl.AuthController.ExecTmp)
	r.Mux.HandleFunc("POST /sign-in", r.Ctl.AuthController.SignIn)
	r.Mux.HandleFunc("GET /sign-out", r.Ctl.AuthController.SignOut)

	// Main page
	// r.Mux.HandleFunc("POST /main", r.Ctl.MainController.)

	// categories routes
	// r.Mux.HandleFunc("GET /categories", r.Ctl.MainController)         // get all categories
	// r.Mux.HandleFunc("GET /categories/{id}", r.Ctl.MainController)    // get a category by id
	// r.Mux.HandleFunc("POST /categories", r.Ctl.MainController)        // create a new category
	// r.Mux.HandleFunc("PUT /categories/{id}", r.Ctl.MainController)    // update a category by id
	// r.Mux.HandleFunc("DELETE /categories/{id}", r.Ctl.MainController) // delete a category by id

	// // posts routes
	// r.Mux.HandleFunc("GET /posts", r.Ctl.MainController)         // get all posts // CRUD
	r.Mux.HandleFunc("GET /posts/{id}", r.Ctl.View)                    // get a post by id
	r.Mux.HandleFunc("GET /create-posts", r.Ctl.CreatePage)            // create a new post
	r.Mux.HandleFunc("POST /posts", r.Ctl.PostController.Create)       // create a new post
	r.Mux.HandleFunc("GET /posts-delete", r.Ctl.PostController.Delete) // delete a post by id
	// Special functi
	// r.Mux.HandleFunc("PUT /posts/{id}", r)    // update a post by id

	// // comments routes
	// r.Mux.HandleFunc("GET /comments", r.Ctl.MainController)         // get all comments
	// r.Mux.HandleFunc("GET /comments/{id}", r.Ctl.MainController)    // get a comment by id
	r.Mux.HandleFunc("POST /comments", r.Ctl.CreateComment) // create a new comment
	// r.Mux.HandleFunc("PUT /comments/{id}", r.Ctl.MainController)    // update a comment by id
	r.Mux.HandleFunc("GET /comments-delete", r.Ctl.DeleteComment) // delete a comment by id

	// Chat
	r.Mux.HandleFunc("GET /chat", r.Ctl.ChatPage)
	// websocket
	r.Mux.HandleFunc("GET /ws/create-pvchat", r.Ctl.WsChatController.CreateRoom)
	// r.Mux.HandleFunc("GET /ws/join-pvchat", r.Ctl.WsChatController.ChatPage)

	// get all chats with last messages by user_id
	// r.Mux.HandleFunc("GET /ws/chats/get-all", r.Ctl.WsChatController.ChatPage)

	// get all online and offline users of the list
	// r.Mux.HandleFunc("GET /ws/chat/users/get-all", r.Ctl.WsChatController.ChatPage)

	// chat op
	// r.Mux.HandleFunc("GET /ws/chat/write-msg", r.Ctl.WsChatController.ChatPage)
	// r.Mux.HandleFunc("GET /ws/chat/read-msg", r.Ctl.WsChatController.ChatPage)
	// r.Mux.HandleFunc("GET /ws/chat/reload-msg", r.Ctl.WsChatController.ChatPage)
	// r.Mux.HandleFunc("GET /ws/chat/more-msgs", r.Ctl.WsChatController.ChatPage)
}

func (r *Router) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// checkAuth
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			slog.Warn(err.Error())
			http.Redirect(w, r, "/sign-in", http.StatusForbidden)
			return
		}
		sessionUUID := cookie.Value

		id, err := GetUserIDFromSession(sessionUUID)
		if err != nil {
			slog.Warn(err.Error())
			http.Redirect(w, r, "/sign-in", http.StatusForbidden)
			return
		}

		// TODO: Make logMiddleware for practice
		// logMiddleware
		slog.Info("User request", "user_id", id, "url", r.URL, "user agent", r.UserAgent(), "user ip", r.RemoteAddr)

		next.ServeHTTP(w, r)
	}
}

func GetUserIDFromSession(sessionUUID string) (int, error) {
	if sessionUUID == "" {
		return 0, fmt.Errorf("session UUID is empty")
	}

	// a new db instance
	dbPath := os.Getenv("DATABASE_FILE_PATH")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	var userId int
	// Выполняем запрос для поиска user_id по sessionUUID
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_uuid=?", sessionUUID).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("Session not found", "sessionUUID", sessionUUID)
			return 0, fmt.Errorf("session not found")
		}
		slog.Error("Error querying session", "error", err)
		return 0, err
	}

	return userId, nil
}
