package controller

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func InitForRelationController() {
	dao.InitDb()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitRelationMQ()
}

// user 4 14 16 17
// 4 -> 14
// 4 -> 16
// 14 -> 16
// 16 -> 4 X
// 17 -> 4
// 4 -> 17 X
// 16 -> 17
// [17 -> 16]
// add check(14 -> 16 16 -> 4 17 -> 16)  find ( 16's follow 17's follower 4's friend)
// delete add check (14 -> 16 16 -> 4 17 -> 16) find (16's follow 17's follower 4's friend)
func TestRelationController(t *testing.T) {
	InitForRelationController()
	e := newExpect(t)
	zywooToken := GetTokenByLogin("zywoo", "7355608", e)
	s1mpleToken := GetTokenByLogin("s1mple", "7355608", e)
	nikoToken := GetTokenByLogin("niko", "8065537", e)
	sh1roToken := GetTokenByLogin("sh1ro", "sh1ro123", e)

	// add
	// 4 -> 14
	resp := e.POST("/douyin/relation/action/").
		WithQuery("token", zywooToken).WithQuery("to_user_id", 14).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// 4 -> 16
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", zywooToken).WithQuery("to_user_id", 16).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// 14 -> 16
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", s1mpleToken).WithQuery("to_user_id", 16).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// 16 -> 4
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", nikoToken).WithQuery("to_user_id", 4).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// 17 -> 4
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", sh1roToken).WithQuery("to_user_id", 4).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// 4 -> 17
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", zywooToken).WithQuery("to_user_id", 17).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// 16 -> 17
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", nikoToken).WithQuery("to_user_id", 17).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// 16's follows are 4 and 17
	userListResp := e.GET("/douyin/relation/follow/list/").
		WithQuery("token", zywooToken).WithQuery("user_id", 16).
		Expect().Status(http.StatusOK).JSON().Object()
	userListResp.Value("status_code").Number().IsEqual(0)
	userListResp.Value("user_list").Array().Length().IsEqual(2)
	userList := userListResp.Value("user_list").Array().Iter()
	for _, u := range userList {
		user := u.Object()
		id := int64(user.Value("id").Number().Raw())
		if id == int64(4) {
			user.Value("name").String().IsEqual("zywoo")
			user.Value("follow_count").Number().IsEqual(3)
			user.Value("follower_count").Number().IsEqual(2)
			user.Value("is_follow").Boolean().IsFalse()
		} else if id == int64(17) {
			user.Value("name").String().IsEqual("sh1ro")
			user.Value("follow_count").Number().IsEqual(1)
			user.Value("follower_count").Number().IsEqual(2)
			user.Value("is_follow").Boolean().IsTrue()
		} else {
			assert.True(t, false)
		}
	}

	// 17's followers are 4 and 16
	userListResp = e.GET("/douyin/relation/follower/list/").
		WithQuery("token", zywooToken).WithQuery("user_id", 17).
		Expect().Status(http.StatusOK).JSON().Object()
	userListResp.Value("status_code").Number().IsEqual(0)
	userListResp.Value("user_list").Array().Length().IsEqual(2)
	userList = userListResp.Value("user_list").Array().Iter()
	for _, u := range userList {
		user := u.Object()
		id := int64(user.Value("id").Number().Raw())
		if id == int64(4) {
			user.Value("name").String().IsEqual("zywoo")
			user.Value("follow_count").Number().IsEqual(3)
			user.Value("follower_count").Number().IsEqual(2)
			user.Value("is_follow").Boolean().IsFalse()
		} else if id == int64(16) {
			user.Value("name").String().IsEqual("niko")
			user.Value("follow_count").Number().IsEqual(2)
			user.Value("follower_count").Number().IsEqual(2)
			user.Value("is_follow").Boolean().IsTrue()
		} else {
			assert.True(t, false)
		}
	}

	// 4's friend are 16 and 17
	userListResp = e.GET("/douyin/relation/friend/list/").
		WithQuery("token", zywooToken).WithQuery("user_id", 4).
		Expect().Status(http.StatusOK).JSON().Object()
	userListResp.Value("status_code").Number().IsEqual(0)
	userListResp.Value("user_list").Array().Length().IsEqual(2)
	userList = userListResp.Value("user_list").Array().Iter()
	for _, u := range userList {
		user := u.Object()
		id := int64(user.Value("id").Number().Raw())
		if id == int64(17) {
			user.Value("name").String().IsEqual("sh1ro")
			user.Value("follow_count").Number().IsEqual(1)
			user.Value("follower_count").Number().IsEqual(2)
			user.Value("is_follow").Boolean().IsTrue()
		} else if id == int64(16) {
			user.Value("name").String().IsEqual("niko")
			user.Value("follow_count").Number().IsEqual(2)
			user.Value("follower_count").Number().IsEqual(2)
			user.Value("is_follow").Boolean().IsTrue()
		} else {
			assert.True(t, false)
		}
	}

	// delete 16 -> 4
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", nikoToken).WithQuery("to_user_id", 4).WithQuery("action_type", 2).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// delete 4 -> 17
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", zywooToken).WithQuery("to_user_id", 17).WithQuery("action_type", 2).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// add 17 -> 16
	resp = e.POST("/douyin/relation/action/").
		WithQuery("token", sh1roToken).WithQuery("to_user_id", 16).WithQuery("action_type", 1).
		Expect().Status(http.StatusOK).JSON().Object()
	resp.Value("status_code").Number().IsEqual(0)

	// 16's follow is 17
	userListResp = e.GET("/douyin/relation/follow/list/").
		WithQuery("token", zywooToken).WithQuery("user_id", 16).
		Expect().Status(http.StatusOK).JSON().Object()
	userListResp.Value("status_code").Number().IsEqual(0)
	userListResp.Value("user_list").Array().Length().IsEqual(1)
	userList = userListResp.Value("user_list").Array().Iter()
	for _, u := range userList {
		user := u.Object()
		id := int64(user.Value("id").Number().Raw())
		if id == int64(17) {
			user.Value("name").String().IsEqual("sh1ro")
			user.Value("follow_count").Number().IsEqual(2)
			user.Value("follower_count").Number().IsEqual(1)
			user.Value("is_follow").Boolean().IsFalse()
		} else {
			assert.True(t, false)
		}
	}

	// 17's follower is 16
	userListResp = e.GET("/douyin/relation/follower/list/").
		WithQuery("token", zywooToken).WithQuery("user_id", 17).
		Expect().Status(http.StatusOK).JSON().Object()
	userListResp.Value("status_code").Number().IsEqual(0)
	userListResp.Value("user_list").Array().Length().IsEqual(1)
	userList = userListResp.Value("user_list").Array().Iter()
	for _, u := range userList {
		user := u.Object()
		id := int64(user.Value("id").Number().Raw())
		if id == int64(16) {
			user.Value("name").String().IsEqual("niko")
			user.Value("follow_count").Number().IsEqual(1)
			user.Value("follower_count").Number().IsEqual(3)
			user.Value("is_follow").Boolean().IsTrue()
		} else {
			assert.True(t, false)
		}
	}

	// 4's friend is none
	userListResp = e.GET("/douyin/relation/friend/list/").
		WithQuery("token", zywooToken).WithQuery("user_id", 4).
		Expect().Status(http.StatusOK).JSON().Object()
	userListResp.Value("status_code").Number().IsEqual(0)
	userListResp.Value("user_list").Array().IsEmpty()
}
