package main

import (
	"douyin/dao"
	"douyin/middleware/ftp"
	"github.com/gin-gonic/gin"
)

func main() {
	//go service.RunMessageServer()
	dao.InitDb()
	ftp.InitFtp()
	ftp.InitSSH()
	r := gin.Default()
	initRouter(r)
	r.Run()
}
