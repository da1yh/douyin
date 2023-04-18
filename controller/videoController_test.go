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

//func TestPublishAction(t *testing.T) {
//	e := newExpect(t)
//	userName, password := "zywoo", "7355608"
//	loginResp := e.POST("/douyin/user/login/").
//		WithQuery("username", userName).WithQuery("password", password).
//		Expect().Status(http.StatusOK).JSON().Object()
//	loginResp.Value("status_code").IsEqual(0)
//	loginResp.Value("user_id").Number().Gt(0)
//	loginResp.Value("token").String().Length().Gt(0)
//	token := loginResp.Value("token").String().Raw()
//	publishResp := e.POST("/douyin/publish/action/").WithMultipart().
//		WithFormField("token", token).WithFormField("title", "iamabear").
//		WithFile("data", "../public/bear.mp4").Expect().Status(http.StatusOK).JSON().Object()
//	publishResp.Value("status_code").IsEqual(0)
//}

func TestPublishList(t *testing.T) {
	e := newExpect(t)
	userName, password := "zywoo", "7355608"
	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", userName).WithQuery("password", password).
		Expect().Status(http.StatusOK).JSON().Object()
	loginResp.Value("status_code").IsEqual(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)
	token := loginResp.Value("token").String().Raw()

	publishListResp := e.GET("/douyin/publish/list/").
		WithQuery("user_id", 16).WithQuery("token", token).
		Expect().StatusList(http.StatusOK).JSON().Object()

	publishListResp.Value("status_code").IsEqual(0)
	publishListResp.Value("video_list").Array().Length().IsEqual(2)
	videoList := publishListResp.Value("video_list").Array().Iter()
	for _, v := range videoList {
		video := v.Object()
		video.Value("title").String().IsEqual("bear")
		video.Value("favorite_count").Number().IsEqual(0)
		video.Value("comment_count").Number().IsEqual(0)
		video.Value("is_favorite").Boolean().IsEqual(false)
		author := video.Value("author").Object()
		author.Value("id").Number().IsEqual(16)
		author.Value("name").String().IsEqual("niko")
	}
}
