package main

import (
	"log"
)

func main() {
	cfg := Config{
		addr: ":8080",
	}

	app := &Application{
		config: cfg,
	}

	mux := app.Mount()
	log.Fatal(app.Run(mux))
}
