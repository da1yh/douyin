package redis

import (
	"context"
	"douyin/config"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RedisCli *redis.Client

func InitRedis() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     config.RedisServerAddr,
		Password: config.RedisServerPwd,
		DB:       0,
	})

	_, err := RedisCli.Ping(Ctx).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("connect to redis server successfully")
	}
}
