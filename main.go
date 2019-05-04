package main

import (
	"flag"
	"log"

	"shier/internal/api"
	"shier/pkg/db"
	"shier/pkg/redis"
)

var (
	listenPort = flag.String("listen-port", "9000", "Port where app listen to")
	dbUrl      = flag.String("db-url", "postgres://docker:docker@localhost:5432/shierdb?sslmode=disable", "Connection string to postgres")
	debug      = flag.Bool("debug", true, "Want to verbose query or not")
	redisAddr  = flag.String("redis-addr", ":6379", "Address string to redis")
	redisPass  = flag.String("redis-pass", "", "Password string to redis")
	redisDb    = flag.Int("redis-db", 0, "DB integer to redis")
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

	// redis configurations
	redisConfig := redis.Config{
		Addr:     *redisAddr,
		Password: *redisPass,
		DB:       *redisDb,
	}

	// redis connection
	err = redis.NewConnection(&redisConfig)
	if err != nil {
		log.Fatalf("redis connection failed")
	}
	defer redis.GetConnection().Close()

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
