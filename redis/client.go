package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var Client *redis.Client

func Init() {
	url := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(url)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	Client = redis.NewClient(opt)
}
