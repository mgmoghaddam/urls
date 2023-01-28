package configs

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

var client *redis.Client

func InitRedis() *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", vi.Get("redis.host"), vi.Get("redis.port")),
		Password: vi.GetString("redis.password"),
		DB:       vi.GetInt("redis.db"),
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis %v", err)
	}
	log.Println(pong, err)

	return client
}

func GetRedisClient() *redis.Client {
	return client
}
