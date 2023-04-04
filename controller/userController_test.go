package controller

import (
	"douyin/dao"
	"net/http"
	"testing"
)

func TestInit(t *testing.T) {
	dao.InitDb()
}

//func TestUserRegister(t *testing.T) {
//	e := newExpect(t)
//	userName, password := "niko", "8065537"
//	registerResp := e.POST("/douyin/user/register/").
//		WithQuery("username", userName).WithQuery("password", password).
//		WithFormField("username", userName).WithFormField("password", password).
//		Expect().Status(http.StatusOK).JSON().Object()
//	registerResp.Value("status_code").Number().IsEqual(0)
//	registerResp.Value("user_id").Number().Gt(0)
//	registerResp.Value("token").String().Length().Gt(0)
//}

func TestUserLogin(t *testing.T) {
	e := newExpect(t)
	userName, password := "zywoo", "7355608"
	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", userName).WithQuery("password", password).
		Expect().Status(http.StatusOK).JSON().Object()
	loginResp.Value("status_code").IsEqual(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)
}
