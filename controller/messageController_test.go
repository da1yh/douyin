package controller

import (
	"douyin/dao"
	"douyin/middleware/redis"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func SLEEP() {
	time.Sleep(3 * time.Second)
}

// 4 -> 14 'hello s1mple, i am zywoo'
// 4 -> 16 'hello niko, i am zywoo'
// 14 -> 4 'hi zywoo, i am s1mple'
// 4 -> 14 'you are 🐐'
func TestMessageController(t *testing.T) {
	dao.InitDb()
	redis.InitRedis()
	e := newExpect(t)
	zywooToken := GetTokenByLogin("zywoo", "7355608", e)
	s1mpleToken := GetTokenByLogin("s1mple", "7355608", e)
	_ = GetTokenByLogin("niko", "8065537", e)
	resp := e.POST("/douyin/message/action/").
		WithQuery("token", zywooToken).WithQuery("to_user_id", 14).
		WithQuery("action_type", 1).WithQuery("content", "hello s1mple, i am zywoo").
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	SLEEP()

	resp = e.POST("/douyin/message/action/").
		WithQuery("token", zywooToken).WithQuery("to_user_id", 16).
		WithQuery("action_type", 1).WithQuery("content", "hello niko, i am zywoo").
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	SLEEP()

	resp = e.POST("/douyin/message/action/").
		WithQuery("token", s1mpleToken).WithQuery("to_user_id", 4).
		WithQuery("action_type", 1).WithQuery("content", "hi zywoo, i am s1mple").
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	SLEEP()

	resp = e.POST("/douyin/message/action/").
		WithQuery("token", zywooToken).WithQuery("to_user_id", 14).
		WithQuery("action_type", 1).WithQuery("content", "you are 🐐").
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// SLEEP()

	chatResp := e.GET("/douyin/message/chat/").
		WithQuery("token", s1mpleToken).WithQuery("to_user_id", 4).WithQuery("pre_msg_time", time.Now().Unix()).
		Expect().Status(http.StatusOK).JSON().Object()
	chatResp.Value("status_code").Number().IsEqual(0)
	chatResp.Value("message_list").Array().Length().IsEqual(3)

	// 这里没有测试评论时间倒序排
	for _, m := range chatResp.Value("message_list").Array().Iter() {
		message := m.Object()
		fromUserId := int64(message.Value("from_user_id").Number().Raw())
		toUserId := int64(message.Value("to_user_id").Number().Raw())
		content := message.Value("content").String().Raw()
		if fromUserId == int64(4) && toUserId == int64(14) {
			if content == "hello s1mple, i am zywoo" || content == "you are 🐐" {

			} else {
				assert.True(t, false)
			}
		} else if fromUserId == int64(14) && toUserId == int64(4) {
			if content == "hi zywoo, i am s1mple" {

			} else {
				assert.True(t, false)
			}
		} else {
			assert.True(t, false)
		}
	}

	// check database and redis
}
