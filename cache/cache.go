package cache

import (
	"fmt"
	"go-image/config"
	"go-image/convert"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client
var maxCache uint
var expireTime uint
var IsCache bool

func init() {

	IsCache = convert.StringToBool(config.GetSetting("redis.cache"))
	maxCache = convert.StringToUint(config.GetSetting("redis.max_cache"))
	expireTime = convert.StringToUint(config.GetSetting("redis.expire_time"))

	if IsCache {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.GetSetting("redis.addr"),
			Password: config.GetSetting("redis.password"),
			DB:       convert.StringToInt(config.GetSetting("redis.db")),
		})

		if _, err := redisClient.Ping().Result(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

}

func Set(key string, value interface{}) {
	if uint(len(value.([]byte))) <= maxCache {
		err := redisClient.Set(key, value, time.Second*time.Duration(expireTime)).Err()
		if err != nil {
			log.Println(err)
		}
	}

}

func Get(key string) *[]byte {
	val, err := redisClient.Get(key).Bytes()
	if err == redis.Nil || err == nil {
		return &val
	}

	return nil
}

func Del(key string) {
	vals, err := redisClient.Keys(key + ":*").Result()
	if err == redis.Nil || err == nil {
		redisClient.Del(vals...)
	}
}
