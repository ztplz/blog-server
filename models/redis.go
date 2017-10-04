package models

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// RedisClient redis连接
var RedisClient *redis.Client

// InitialRedis   初始化连接redis
func InitialRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		// redis 地址
		Addr: "localhost:6379",

		// redis 连接密码
		// Password: "123456",

		// redis 连接的数据库,使用默认的数据库
		DB: 0,
	})

	pong, err := RedisClient.Ping().Result()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Redis ping failed")
	}

	log.WithFields(log.Fields{
		"message": pong,
	}).Info("Redis connect success")
}
