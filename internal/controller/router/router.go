package router

import (
	"net/http"

	"github.com/Pruel/real-time-forum/internal/controller"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

type Router struct { // public, protected, private -> Router
	Mux *http.ServeMux
	Ctl *controller.Controller
	// Chi *chi.Router
	// gorrilaMux *gorilla.Router
}

// New create and return a new instance of the ServeMux router
func New(db *sqlite.Database) *Router {
	return &Router{
		Mux: http.NewServeMux(),
		Ctl: controller.New(db),
	}
}

func (r *Router) InitRouter() {

	// statics
	r.Mux.HandleFunc("GET /", r.Ctl.MainController)

	// auth routes, sign-up, sign-in, sign-out
	r.Mux.HandleFunc("GET /sign-up", r.Ctl.AuthController.SignUpPage)
	r.Mux.HandleFunc("POST /sign-up", r.Ctl.AuthController.SignUp)
	r.Mux.HandleFunc("/sign-in", r.Ctl.MainController)
	r.Mux.HandleFunc("/sign-out", r.Ctl.MainController)

	// categories routes
	r.Mux.HandleFunc("GET /categories", r.Ctl.MainController)         // get all categories
	r.Mux.HandleFunc("GET /categories/{id}", r.Ctl.MainController)    // get a category by id
	r.Mux.HandleFunc("POST /categories", r.Ctl.MainController)        // create a new category
	r.Mux.HandleFunc("PUT /categories/{id}", r.Ctl.MainController)    // update a category by id
	r.Mux.HandleFunc("DELETE /categories/{id}", r.Ctl.MainController) // delete a category by id

	// posts routes
	r.Mux.HandleFunc("GET /posts", r.Ctl.MainController)         // get all posts // CRUD
	r.Mux.HandleFunc("GET /posts/{id}", r.Ctl.MainController)    // get a post by id
	r.Mux.HandleFunc("POST /posts", r.Ctl.MainController)        // create a new post
	r.Mux.HandleFunc("PUT /posts/{id}", r.Ctl.MainController)    // update a post by id
	r.Mux.HandleFunc("DELETE /posts/{id}", r.Ctl.MainController) // delete a post by id

	// comments routes
	r.Mux.HandleFunc("GET /comments", r.Ctl.MainController)         // get all comments
	r.Mux.HandleFunc("GET /comments/{id}", r.Ctl.MainController)    // get a comment by id
	r.Mux.HandleFunc("POST /comments", r.Ctl.MainController)        // create a new comment
	r.Mux.HandleFunc("PUT /comments/{id}", r.Ctl.MainController)    // update a comment by id
	r.Mux.HandleFunc("DELETE /comments/{id}", r.Ctl.MainController) // delete a comment by id
}
