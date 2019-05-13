package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
)

type Server struct {
	DB     *pg.DB
	Router *gin.Engine
	Redis  *redis.Client
}

func (s *Server) InitDB(url string, debug bool) error {
	opt, err := pg.ParseURL(url)
	if err != nil {
		return err
	}

	s.DB = pg.Connect(opt)

	// ensure database connection is successfully
	_, err = s.DB.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) InitRedis(addr, password string, db int) error {
	s.Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := s.Redis.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}
