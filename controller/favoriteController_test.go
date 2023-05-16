package controller

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
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

func TestFavoriteList(t *testing.T) {
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

	// 验证是否查询空的点赞列表会出错
	favoriteListResp := e.GET("/douyin/favorite/list/").
		WithQuery("user_id", 14).WithQuery("token", token).
		Expect().Status(http.StatusOK).JSON().Object()
	favoriteListResp.Value("status_code").Number().IsEqual(0)
	favoriteListResp.Value("video_list").Array().Length().IsEqual(0)

	// 模拟用户登录以点赞，s1mple点赞1视频和2视频，我（zywoo）点赞1视频
	userName, password = "s1mple", "7355608"
	loginResp = e.POST("/douyin/user/login/").
		WithQuery("username", userName).WithQuery("password", password).
		Expect().Status(http.StatusOK).JSON().Object()
	loginResp.Value("status_code").Number().IsEqual(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)
	anotherToken := loginResp.Value("token").String().Raw()

	resp := e.POST("/douyin/favorite/action/").
		WithQuery("token", anotherToken).WithQuery("video_id", 1).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)
	resp = e.POST("/douyin/favorite/action/").
		WithQuery("token", anotherToken).WithQuery("video_id", 2).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)
	resp = e.POST("/douyin/favorite/action/").
		WithQuery("token", token).WithQuery("video_id", 1).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	//=========================================================

	favoriteListResp = e.GET("/douyin/favorite/list/").
		WithQuery("user_id", 14).WithQuery("token", token).
		Expect().Status(http.StatusOK).JSON().Object()
	favoriteListResp.Value("status_code").Number().IsEqual(0)
	favoriteListResp.Value("video_list").Array().Length().IsEqual(2)
	favoriteList := favoriteListResp.Value("video_list").Array().Iter()
	one, two := 0, 0
	for _, favorite := range favoriteList {
		video := favorite.Object()
		id := int64(video.Value("id").Number().Raw())
		if id == int64(1) {
			one++
			video.Value("favorite_count").Number().IsEqual(2)
			video.Value("is_favorite").Boolean().IsEqual(true)
			video.Value("title").String().IsEqual("bear")
			video.Value("author").Object().Value("id").Number().IsEqual(14)
		} else {
			two++
			video.Value("favorite_count").Number().IsEqual(1)
			video.Value("is_favorite").Boolean().IsEqual(false)
			video.Value("title").String().IsEqual("bear")
			video.Value("author").Object().Value("id").Number().IsEqual(14)
		}
	}
	assert.Equal(t, one, 1)
	assert.Equal(t, two, 1)

	// 模拟s1mple取消点赞
	resp = e.POST("/douyin/favorite/action/").
		WithQuery("token", anotherToken).WithQuery("video_id", 1).WithQuery("action_type", 2).
		Expect().Status(http.StatusOK).JSON().Object()
	resp = e.POST("/douyin/favorite/action/").
		WithQuery("token", anotherToken).WithQuery("video_id", 2).WithQuery("action_type", 2).
		Expect().Status(http.StatusOK).JSON().Object()
	favoriteListResp = e.GET("/douyin/favorite/list/").
		WithQuery("user_id", 14).WithQuery("token", token).
		Expect().Status(http.StatusOK).JSON().Object()
	favoriteListResp.Value("status_code").Number().IsEqual(0)
	favoriteListResp.Value("video_list").Array().IsEmpty()
}

func TestGetVideoById(t *testing.T) {
	Init()
	videoList := make([]Video, 0)
	var wg sync.WaitGroup
	wg.Add(1)
	GetVideoById(1, &videoList, &wg, 4)
	assert.Equal(t, len(videoList), 1)
	video := videoList[0]
	assert.Equal(t, video.Title, "bear")
	assert.Equal(t, video.Id, int64(1))
	assert.True(t, video.IsFavorite)
	assert.Equal(t, video.FavoriteCount, int64(2))
	assert.Equal(t, video.CommentCount, int64(0))
	//assert.Equal(t, video.Author.Id, int64(16))
}

func TestGetVideoListByIds(t *testing.T) {
	Init()
	ids := []int64{1, 2}
	videoList := GetVideoListByIds(ids, 4)
	assert.Equal(t, len(videoList), 2)
	one, two := 0, 0
	for _, video := range videoList {
		if video.Id == int64(1) {
			one++
			assert.Equal(t, video.FavoriteCount, int64(2))
			assert.True(t, video.IsFavorite)
		} else if video.Id == int64(2) {
			two++
			assert.Equal(t, video.FavoriteCount, int64(1))
			assert.False(t, video.IsFavorite)
		}
	}
	assert.Equal(t, one, 1)
	assert.Equal(t, two, 1)
}
