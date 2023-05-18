package main

import (
	"douyin/dao"
	"douyin/middleware/ftp"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	//go service.RunMessageServer()
	dao.InitDb()
	ftp.InitFtp()
	ftp.InitSSH()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFavoriteMQ()
	rabbitmq.InitCommentMQ()
	r := gin.Default()
	initRouter(r)
	time.Sleep(3 * time.Second)
	r.Run()
}
