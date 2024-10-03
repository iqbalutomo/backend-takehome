package config

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedistClient() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:       opt.Addr,
		Password:   opt.Password,
		DB:         opt.DB,
		MaxRetries: 3,
	})

	return client
}
