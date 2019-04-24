package main

import (
	"flag"
	"log"

	"github.com/paramahastha/shier/internal/api"
	"github.com/paramahastha/shier/pkg/db"
)

var (
	listenPort = flag.String("listen-port", "9000", "Port where app listen to")
	dbUrl      = flag.String("db-url", "postgres://docker:docker@localhost:5432/shierdb?sslmode=disable", "Connection string to postgres")
	debug      = flag.Bool("debug", true, "Want to verbose query or not")
)

func main() {
	flag.Parse()

	// database configurations
	dbConfig := db.Config{
		URL:   *dbUrl,
		Debug: *debug,
	}

	// database connection
	dbConn, err := db.NewConnection(&dbConfig)
	if err != nil {
		log.Fatalf("database connection failed")
	}
	defer dbConn.Close()

	// run migrations
	err = db.Migrate(*dbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// server configurations
	apiConfig := api.Config{
		ListenPort: *listenPort,
	}

	apiConfig.Start() // start server
}
