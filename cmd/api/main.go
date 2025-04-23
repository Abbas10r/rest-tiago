package main

import (
	"log"
	"os"
	"socialApp/internal/auth"
	"socialApp/internal/mailer"
	"socialApp/internal/store"
	"socialApp/internal/store/db"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const version = "0.0.2"

// @title Social API
// @description API for Social API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath	/
//
// @securityDefinitions.apikey  ApiKeyAuth
// @in 				header
// @name 			Authorization
// @description
func main() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := Config{
		addr:        ":8080",
		frontendURL: os.Getenv("FRONTEND_URL"),
		db: dbConfig{
			addr:         os.Getenv("DB_CONN"),
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
		env: "production",
		mail: mailConfig{
			exp: time.Hour * 24 * 3, //3 days
			sendGrid: sendGridConfig{
				apiKey: os.Getenv("SENDGRID_API_KEY"),
			},
			fromEmail: os.Getenv("FROM_EMAIL"),
		},
		auth: authConfig{
			token: tokenConfig{
				secret: "secret",
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "gophersocial",
			},
		},
	}

	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(db)
	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	JWTAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &Application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: JWTAuthenticator,
	}

	mux := app.Mount()
	logger.Fatal(app.Run(mux))
}
