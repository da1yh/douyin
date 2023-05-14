package controller

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"net/http"
	"testing"
)

func Init() {
	dao.InitDb()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFavoriteMQ()
}

func TestFavoriteAction(t *testing.T) {
	Init()
	e := newExpect(t)
	userName, password := "zywoo", "7355608"
	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", userName).WithQuery("password", password).
		Expect().Status(http.StatusOK).JSON().Object()
	loginResp.Value("status_code").Number().IsEqual(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)
	token := loginResp.Value("token").String().Raw()
	resp := e.POST("/douyin/favorite/action/").
		WithQuery("token", token).WithQuery("video_id", 3).WithQuery("action_type", 2).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

}
