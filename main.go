package main

import (
	"flag"
	"log"

	"shier/internal/api"
	"shier/pkg/db"
)

var (
	listenPort = flag.String("listen-port", "9000", "Port where app listen to")
	dbUrl      = flag.String("db-url", "postgres://docker:docker@localhost:5432/shierdb?sslmode=disable", "Connection string to postgres")
	debug      = flag.Bool("debug", true, "Want to verbose query or not")
	redisAddr  = flag.String("redis-addr", ":6000", "Address string to redis")
	redisPass  = flag.String("redis-pass", "", "Password string to redis")
	redisDb    = flag.Int("redis-db", 0, "DB integer to redis")
)

func main() {
	flag.Parse()

	// Init server
	server := api.Server{}

	// database connection
	err := server.InitDB(*dbUrl, *debug)
	if err != nil {
		log.Fatalf("database connection failed")
	}
	defer server.DB.Close()

	// redis connection
	err = server.InitRedis(*redisAddr, *redisPass, *redisDb)
	if err != nil {
		log.Fatalf("redis connection failed")
	}
	defer server.Redis.Close()

	// run migrations
	err = db.Migrate(*dbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}

	server.SetupRouter(*listenPort) // start server
}
