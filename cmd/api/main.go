package main

import (
	"log"
	"os"
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
// @BasePath
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
		addr: ":8080",
		db: dbConfig{
			addr:         os.Getenv("DB_CONN"),
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
		env: "Development",
		mail: mailConfig{
			exp: time.Hour * 24 * 3, //3 days
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

	app := &Application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.Mount()
	logger.Fatal(app.Run(mux))
}
