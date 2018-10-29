package router

import (
	. "common/http_helpers"
	"web/handler"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	// init router
	r := mux.NewRouter()
	// rule helper
	add := func(method string, path string, f HandlerFunc) {
		r.Methods(method).Path(path).HandlerFunc(Wrap(f))
	}
	addTHF := func(method string, path string, f ThinHandlerFunc) {
		add(method, path, WrapTHF(f))
	}

	//FAKE
	addTHF("POST", "/user/{name}", handler.CreateUser)

	return r
}
