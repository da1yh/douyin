package test

import (
	"github.com/stretchr/testify/assert"
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

func TestComment(t *testing.T) {
	// get videos by feeding and comment the first video ->
	// . create a user using getTestUser()
	// . comment the video, check the status_code and comment_id of response
	// get the comment list of corresponding video, check the status_code and related keys
	//check if new comment is in comment list
	// delete corresponding comment, check the status_code
	e := newExpect(t)
	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().IsEqual(0)
	feedResp.Value("video_list").Array().Length().Gt(0)
	firstVideo := feedResp.Value("video_list").Array().Value(0).Object()
	videoId := firstVideo.Value("id").Number().Raw()

	_, token := getTestUser("ropz", e)
	commentResp := e.POST("/douyin/comment/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 1).WithQuery("comment_text", "cs2 is on the way").
		Expect().Status(http.StatusOK).JSON().Object()
	commentResp.Value("status_code").Number().IsEqual(0)
	commentResp.Value("comment").Object().Value("id").Number().Gt(0)
	commentId := int(commentResp.Value("comment").Object().Value("id").Number().Raw())

	commentListResp := e.GET("/douyin/comment/list/").
		WithQuery("token", token).WithQuery("video_id", videoId).
		Expect().Status(http.StatusOK).JSON().Object()
	commentListResp.Value("status_code").Number().IsEqual(0)
	hasComment := false
	for _, element := range commentListResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		comment.ContainsKey("id")
		comment.ContainsKey("user")
		comment.Value("content").String().NotEmpty()
		comment.Value("create_date").String().NotEmpty()
		if int(comment.Value("id").Number().Raw()) == commentId {
			hasComment = true
		}
	}
	assert.True(t, hasComment, "can't find test comment")
	delCommentResp := e.POST("/douyin/comment/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 2).WithQuery("comment_id", commentId).
		Expect().Status(http.StatusOK).JSON().Object()
	delCommentResp.Value("status_code").Number().IsEqual(0)
}
