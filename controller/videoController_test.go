package controller

import (
	"douyin/dao"
	"douyin/middleware/ftp"
	"net/http"
	"testing"
)

func TestInit2(t *testing.T) {
	dao.InitDb()
	ftp.InitFtp()
	ftp.InitSSH()
}

func TestPublishAction(t *testing.T) {
	e := newExpect(t)
	userName, password := "zywoo", "7355608"
	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", userName).WithQuery("password", password).
		Expect().Status(http.StatusOK).JSON().Object()
	loginResp.Value("status_code").IsEqual(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)
	token := loginResp.Value("token").String().Raw()
	publishResp := e.POST("/douyin/publish/action/").WithMultipart().
		WithFormField("token", token).WithFormField("title", "iamabear").
		WithFile("data", "../public/bear.mp4").Expect().Status(http.StatusOK).JSON().Object()
	publishResp.Value("status_code").IsEqual(0)
}
