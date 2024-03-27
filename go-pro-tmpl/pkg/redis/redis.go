package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

func (rc *RedisClient) Get(key string) (string, error) {
	val, err := rc.Client.Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("redis key not found")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (rc *RedisClient) Set(key, value string) error {
	return rc.Client.Set(key, value, 0).Err()
}

func (rc *RedisClient) KeyExists(key string) (bool, error) {
	val, err := rc.Client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

func (rc *RedisClient) UpdateValue(key, value string) error {
	return rc.Client.Set(key, value, 0).Err()
}

func (rc *RedisClient) DeleteKey(key string) (int64, error) {
	val, err := rc.Client.Del(key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}
