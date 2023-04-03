package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	//go service.RunMessageServer()
	r := gin.Default()
	initRouter(r)
	r.Run()
}
