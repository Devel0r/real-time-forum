package router

import (
	"fmt"
	"github.com/Pruel/real-time-forum/internal/controller"
	"net/http"
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

	// statics
	// r.Mux.HandleFunc("GET /", r.Ctl.MainController) // template -> router -> controller -> model  -> repository -> database

	// auth routes, sign-up, sign-in, sign-out
	r.Mux.HandleFunc("/api/signup", r.Ctl.AuthController.SignUp) // POST - SignUpPage
	r.Mux.HandleFunc("/api/login", r.Ctl.AuthController.SignIn)
	r.Mux.HandleFunc("/api/logout", r.Ctl.SignOut)
	r.Mux.HandleFunc("/api/check-auth", r.Ctl.AuthController.CheckAuth)

	staticPath := r.Ctl.GetStaticPath()
	fs := http.FileServer(http.Dir(staticPath))
	r.Mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("Serving static file: %s\n", req.URL.Path)
		w.Header().Add("X-Frame-Options", "DENY")
		fs.ServeHTTP(w, req)
	})))

	// Обслуживание SPA
	r.Mux.HandleFunc("/", r.Ctl.MainController)
	// categories routes // roles, categories, items, -> ?
	// r.Mux.HandleFunc("GET /categories", r.Ctl.MainController)         // get all categories
	// r.Mux.HandleFunc("GET /categories/{id}", r.Ctl.MainController)    // get a category by id
	// r.Mux.HandleFunc("POST /categories", r.Ctl.MainController)        // create a new category
	// r.Mux.HandleFunc("PUT /categories/{id}", r.Ctl.MainController)    // update a category by id
	// r.Mux.HandleFunc("DELETE /categories/{id}", r.Ctl.MainController) // delete a category by id

	// Час вышел, так что всё ок чпок, главное пофикси эту штуку а то я умру, она у меня всю душу высосет
	// У тебя помойму трафик отлетел
	// I did write you in telegram
	// posts routes
	r.Mux.HandleFunc("GET /posts", r.Ctl.MainController)         // get all posts // CRUD
	r.Mux.HandleFunc("GET /posts/{id}", r.Ctl.MainController)    // get a post by id
	r.Mux.HandleFunc("POST /posts", r.Ctl.MainController)        // create a new post
	r.Mux.HandleFunc("PUT /posts/{id}", r.Ctl.MainController)    // update a post by id
	r.Mux.HandleFunc("DELETE /posts/{id}", r.Ctl.MainController) // delete a post by id

	// // comments routes
	// r.Mux.HandleFunc("GET /comments", r.Ctl.MainController)         // get all comments
	// r.Mux.HandleFunc("GET /comments/{id}", r.Ctl.MainController)    // get a comment by id
	// r.Mux.HandleFunc("POST /comments", r.Ctl.MainController)        // create a new comment
	// r.Mux.HandleFunc("PUT /comments/{id}", r.Ctl.MainController)    // update a comment by id
	// r.Mux.HandleFunc("DELETE /comments/{id}", r.Ctl.MainController) // delete a comment by id
}
