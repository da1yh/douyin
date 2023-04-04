package main

import (
	"douyin/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	//go service.RunMessageServer()
	dao.InitDb()
	r := gin.Default()
	initRouter(r)
	r.Run()
}
