package main

import (
	"fmt"
	"net/http"
	"socialApp/cmd/api/docs"
	"socialApp/internal/auth"
	"socialApp/internal/mailer"
	"socialApp/internal/store"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Application struct {
	config        Config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type Config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
}

type authConfig struct {
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	exp       time.Duration
	fromEmail string
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *Application) Mount() *mux.Router {
	mux := mux.NewRouter()

	docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
	mux.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(docsURL))).Methods(http.MethodGet)

	healthRouter := mux.PathPrefix("/v1/health").Subrouter()
	healthRouter.Use(app.AuthTokenAuthMiddleware)
	healthRouter.HandleFunc("", app.HealthCheckHandler).Methods("GET") //curl http://localhost:8080/v1/health

	p := mux.PathPrefix("/v1/posts").Subrouter()
	p.Use(app.AuthTokenAuthMiddleware)
	p.HandleFunc("", app.createPostHandler).Methods("POST")
	p2 := p.PathPrefix("").Subrouter()
	p2.Use(app.postsContextMiddleware)
	p2.HandleFunc("/{id}", app.getPostHandler).Methods("GET")

	p3 := p2.PathPrefix("/{id}").Subrouter()
	deleteHandler := app.checkPostOwnership("admin", http.HandlerFunc(app.deletePostHandler))
	p3.Handle("", deleteHandler).Methods("DELETE")
	updateHandler := app.checkPostOwnership("moderator", http.HandlerFunc(app.updatePostHandler))
	p3.Handle("", updateHandler).Methods("PUT")

	u := mux.PathPrefix("/v1/users").Subrouter()
	u.Use(app.AuthTokenAuthMiddleware)
	u.HandleFunc("/{id}", app.getUserHandler).Methods("GET")

	f := mux.PathPrefix("/v1/users/{id}").Subrouter()
	f.Use(app.AuthTokenAuthMiddleware)
	f.HandleFunc("/follow", app.followUserHandler).Methods("PUT")
	f.HandleFunc("/unfollow", app.unfollowUserHandler).Methods("PUT")

	feed := mux.PathPrefix("/v1/feed").Subrouter()
	feed.Use(app.AuthTokenAuthMiddleware)
	feed.HandleFunc("", app.getUserFeedHandler).Methods("GET")

	mux.HandleFunc("/users/activate/{token}", app.activateUserHandler).Methods("PUT")

	authRouter := mux.PathPrefix("/authentication").Subrouter()
	authRouter.HandleFunc("/user", app.registerUserHandler).Methods("POST")
	authRouter.HandleFunc("/token", app.createTokenHandler).Methods("POST")
	mux.NewRoute()
	return mux
}

func (app *Application) Run(mux *mux.Router) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has started at ", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}
