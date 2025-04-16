package main

import (
	"log"
	"net/http"
	"socialApp/internal/store"
	"time"

	"github.com/gorilla/mux"
)

type Application struct {
	config Config
	store  store.Storage
}

type Config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *Application) Mount() *mux.Router {
	mux := mux.NewRouter()

	mux.HandleFunc("/v1/health", app.HealthCheckHandler).Methods("GET") //curl http://localhost:8080/v1/health

	p := mux.PathPrefix("/v1/posts").Subrouter()
	p.HandleFunc("", app.createPostHandler).Methods("POST")
	p.HandleFunc("/{id}", app.getPostHandler).Methods("GET")
	p.HandleFunc("/{id}", app.deletePostHandler).Methods("DELETE")
	p.HandleFunc("/{id}", app.updatePostHandler).Methods("PUT")

	u := mux.PathPrefix("/v1/users").Subrouter()
	u.HandleFunc("/{id}", app.getUserHandler).Methods("GET")
	mux.NewRoute()
	return mux
}

func (app *Application) Run(mux *mux.Router) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
