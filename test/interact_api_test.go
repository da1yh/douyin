package test

import (
	"net/http"
	"testing"
)

func TestFavorite(t *testing.T) {
	// get videos by feeding and thumb up the first video ->
	// . create a user using getTestUser()
	// . thumb up the corresponding video, check the status code of response
	// for favorite list, get the corresponding video list, check the status code of response
	// for every video in video list, check related keys
	e := newExpect(t)
	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().IsEqual(0)
	feedResp.Value("video_list").Array().Length().Gt(0)
	firstVideo := feedResp.Value("video_list").Array().Value(0).Object()
	videoId := firstVideo.Value("id").Number().Raw()

	userId, token := getTestUser("ropz", e)

	favoriteResp := e.POST("/douyin/favorite/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()

	favoriteResp.Value("status_code").Number().IsEqual(0)

	favoriteListResp := e.GET("/douyin/favorite/list/").
		WithQuery("user_id", userId).WithQuery("token", token).
		Expect().Status(http.StatusOK).JSON().Object()
	favoriteListResp.Value("status_code").Number().IsEqual(0)
	for _, element := range favoriteListResp.Value("video_list").Array().Iter() {
		video := element.Object()
		video.ContainsKey("id")
		video.ContainsKey("author")
		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
	}
}
