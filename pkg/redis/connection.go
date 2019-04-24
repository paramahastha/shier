package redis

import "github.com/go-redis/redis"

var client *redis.Client

func NewConnection(c *Config) error {
	client = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

func GetConnection() *redis.Client {
	return client
}
