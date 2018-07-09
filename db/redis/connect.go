package redis

import (
	"github.com/go-redis/redis"
	"../../config"
	"time"
)

var redisCacheClient *redis.Client
var redisDirectoryClient *redis.Client

func init() {
	redisCacheClient = connectCacheRedis()
	redisDirectoryClient = connectDirectoryRedis()
	return
}

func connectCacheRedis() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.GetString("redis.cache.addr"),
		Password: "",
		DB:       config.AppConfig.GetInt("redis.cache.db"),
	})
	return
}

func connectDirectoryRedis() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.GetString("redis.directory.addr"),
		Password: "",
		DB:       config.AppConfig.GetInt("redis.directory.db"),
	})
	return
}

func RedisCacheManager() *redis.Client {
	return redisCacheClient
}

func RedisDirectoryManager() *redis.Client {
	return redisCacheClient
}

func CloseRedis(client *redis.Client) {
	err := client.Close()
	if err != nil {
		panic(err)
	}
	return
}

func GetValue(redisClient *redis.Client, key string) (interface{}, error) {
	var (
		val interface{}
		err error
	)

	if val, err = redisClient.Get(key).Result(); err != nil {
		return val, err
	}

	return val, nil
}

func SetValue(redisClient *redis.Client, key string, value interface{}) error {
	var err error
	if err = redisClient.Set(key, value, 0).Err(); err != nil {
		return err
	}

	return nil
}

func SetValueExpire(redisClient *redis.Client, key string, value interface{}, ex time.Duration) error {
	var err error
	if err = redisClient.Set(key, value, ex).Err(); err != nil {
		return err
	}

	return nil
}
func DelKey(redisClient *redis.Client, key string) error {
	var err error
	if _, err = redisClient.Del(key).Result(); err != nil {
		return err
	}
	return nil
}
