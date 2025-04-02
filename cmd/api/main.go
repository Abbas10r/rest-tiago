package main

import (
	"log"
	"socialApp/internal/store"
	"socialApp/internal/store/db"

	_ "github.com/lib/pq"
)

func main() {
	cfg := Config{
		addr: ":8080",
		db: dbConfig{
			addr:         "postgres://social:social@localhost/social?sslmode=disable",
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Printf("database connection pool established")

	store := store.NewStorage(db)

	app := &Application{
		config: cfg,
		store:  store,
	}

	mux := app.Mount()
	log.Fatal(app.Run(mux))
}
