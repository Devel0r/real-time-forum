package router

import "net/http"

type Router struct { // public, protected, private -> Router
	Mux *http.ServeMux
}

func InitRouter() {

}
