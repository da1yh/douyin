package controller

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"net/http"
	"testing"
)

func CommentInit() {
	dao.InitDb()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitCommentMQ()
}

// 测试评论
// zywoo->1 "good one" X
// zywoo->2 "haha funny"
// zywoo->1 "comment again"
func TestCommentAction(t *testing.T) {
	CommentInit()
	e := newExpect(t)
	userName, password := "zywoo", "7355608"
	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", userName).WithQuery("password", password).
		Expect().Status(http.StatusOK).JSON().Object()
	loginResp.Value("status_code").Number().IsEqual(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)
	token := loginResp.Value("token").String().Raw()

	//commentActionResp := e.POST("/douyin/comment/action/").
	//	WithQuery("token", token).WithQuery("video_id", 1).
	//	WithQuery("action_type", 1).WithQuery("comment_text", "good one").
	//	Expect().Status(http.StatusOK).JSON().Object()
	//commentActionResp.Value("status_code").Number().IsEqual(0)
	//comment := commentActionResp.Value("comment").Object()
	//user := comment.Value("user").Object()
	//comment.Value("content").String().IsEqual("good one")
	//user.Value("id").Number().IsEqual(int64(4))
	////id := comment.Value("id").Number().Raw()
	//
	//commentActionResp = e.POST("/douyin/comment/action/").
	//	WithQuery("token", token).WithQuery("video_id", 2).
	//	WithQuery("action_type", 1).WithQuery("comment_text", "haha funny").
	//	Expect().Status(http.StatusOK).JSON().Object()
	//commentActionResp.Value("status_code").Number().IsEqual(0)
	//comment = commentActionResp.Value("comment").Object()
	//user = comment.Value("user").Object()
	//comment.Value("content").String().IsEqual("haha funny")
	//user.Value("id").Number().IsEqual(int64(4))
	//
	//commentActionResp = e.POST("/douyin/comment/action/").
	//	WithQuery("token", token).WithQuery("video_id", 1).
	//	WithQuery("action_type", 1).WithQuery("comment_text", "comment again").
	//	Expect().Status(http.StatusOK).JSON().Object()
	//commentActionResp.Value("status_code").Number().IsEqual(0)
	//comment = commentActionResp.Value("comment").Object()
	//user = comment.Value("user").Object()
	//comment.Value("content").String().IsEqual("comment again")
	//user.Value("id").Number().IsEqual(int64(4))

	commentActionResp := e.POST("/douyin/comment/action/").
		WithQuery("token", token).WithQuery("video_id", 1).
		WithQuery("action_type", 2).WithQuery("comment_id", 80).
		Expect().Status(http.StatusOK).JSON().Object()
	commentActionResp.Value("status_code").Number().IsEqual(0)

}
