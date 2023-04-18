package controller

import (
	"douyin/dao"
	"douyin/middleware/ftp"
	"net/http"
	"testing"
	"time"
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

func TestFeed(t *testing.T) {
	// test login
	e := newExpect(t)
	userName, password := "zywoo", "7355608"
	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", userName).WithQuery("password", password).
		Expect().Status(http.StatusOK).JSON().Object()
	loginResp.Value("status_code").Number().IsEqual(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)
	token := loginResp.Value("token").String().Raw()

	feedResp := e.GET("/douyin/feed/").
		WithQuery("token", token).WithQuery("latest_time", time.Now().Unix()).
		Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().IsEqual(0)
	feedResp.Value("video_list").Array().Length().IsEqual(3)
	//tmpTime, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2023-04-17 12:06:25"), time.Local)
	//feedResp.Value("next_time").Number().IsEqual(tmpTime.Unix())

	// test visitor

	feedResp = e.GET("/douyin/feed/").
		//WithQuery("token", token).WithQuery("latest_time", time.Now().Unix()).
		Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().IsEqual(0)
	feedResp.Value("video_list").Array().Length().IsEqual(3)
	//tmpTime, _ = time.Parse("2023-04-17 12:06:25", "2023-04-17 12:06:25")
	//feedResp.Value("next_time").Number().IsEqual(tmpTime.Unix())

}
