package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RestHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type RestRouter struct {
	*mux.Router
}

func NewRestRouter() RestRouter {
	return RestRouter{mux.NewRouter()}
}

func (router RestRouter) Register(endpoint string, handler RestHandler) {
	subrouter := router.PathPrefix("/" + endpoint).Subrouter()
	subrouter.Path("").Methods(http.MethodGet).HandlerFunc(handler.Get)
	subrouter.Path("").Methods(http.MethodPost).HandlerFunc(handler.Post)
	subrouter.Path("/{id}").Methods(http.MethodGet).HandlerFunc(handler.GetById)
	subrouter.Path("/{id}").Methods(http.MethodPut).HandlerFunc(handler.Put)
	subrouter.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(handler.Delete)
}
