package redis

import (
	"context"
	"fmt"
	"os"

	redisV8 "github.com/go-redis/redis/v8"
)

// RedisClient = redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_URL")})

// RedisClient is the redis client
var RedisClient *redisV8.Client

// InitRedisClient initializes the redis client
func InitRedisClient() {
	opt, err := redisV8.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	RedisClient = redisV8.NewClient(opt)
}

func ListenRenderInRedis(callback func(string)) {
	pubsub := RedisClient.Subscribe(context.Background(), "render")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}

		fmt.Println(msg.Payload)
		callback(msg.Payload)
	}
}
