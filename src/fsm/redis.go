package fsm

import (
	"strconv"

	"github.com/electric_bayan/weather_bot/config"
	"github.com/redis/go-redis/v9"
)

func New() redis.Client {
	return getRedisClient()
}

func getRedisClient() redis.Client {
	conf := config.New()
	redis_addr := conf.Redis.RedisHost + ":" + strconv.Itoa(conf.Redis.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: conf.Redis.RedisPass,
		DB:       0,
	})
	return *client
}
