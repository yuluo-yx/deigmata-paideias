package redis

import (
	"api-registry/pkg/log"
	"api-registry/pkg/model"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func InitRedis() {

	config := model.GetConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     config.System.Redis.Host + ":" + config.System.Redis.Port,
		Password: config.System.Redis.Password,
		DB:       config.System.Redis.DB,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.SysLog.Warn("Redis connection failed: " + err.Error())
		return
	}
	log.SysLog.Info("Redis connection status: " + pong)

	redisClient = client
}

func GetRedisClient() *redis.Client {

	return redisClient
}
