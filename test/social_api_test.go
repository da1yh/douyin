package test

import (
	"net/http"
	"testing"
)

//func TestRelation(t *testing.T) {
//	// test relation action, get two test users, check the status_code
//	e := newExpect(t)
//	userIdA, tokenA := getTestUser("people", e)
//	userIdB, tokenB := getTestUser("fywoo", e)
//	relationResp := e.POST("/douyin/relation/action/").
//		WithQuery("token", tokenA).WithQuery("to_user_id", userIdB).WithQuery("action_type", 1).
//		Expect().Status(http.StatusOK).JSON().Object()
//	relationResp.Value("status_code").Number().IsEqual(0)
//
//	followResp := e.GET("/douyin/relation/follow/list/").
//		WithQuery("user_id", userIdA).WithQuery("token", tokenA).
//		Expect().Status(http.StatusOK).JSON().Object()
//	followResp.Value("status_code").Number().IsEqual(0)
//	flag := false
//	for _, element := range followResp.Value("user_list").Array().Iter() {
//		user := element.Object()
//		user.ContainsKey("id")
//		if int(user.Value("id").Number().Raw()) == userIdB {
//			flag = true
//		}
//	}
//	assert.True(t, flag, "follow test failed")
//
//	followerResp := e.GET("/douyin/relation/follower/list/").
//		WithQuery("user_id", userIdB).WithQuery("token", tokenB).
//		Expect().Status(http.StatusOK).JSON().Object()
//	followerResp.Value("status_code").Number().IsEqual(0)
//	flag = false
//	for _, element := range followerResp.Value("user_list").Array().Iter() {
//		user := element.Object()
//		user.ContainsKey("id")
//		if int(user.Value("id").Number().Raw()) == userIdA {
//			flag = true
//		}
//	}
//	assert.True(t, flag, "follower test failed")
//}

func TestMessage(t *testing.T) {
	e := newExpect(t)
	userIdA, tokenA := getTestUser("people", e)
	userIdB, tokenB := getTestUser("fywoo", e)
	messageResp := e.POST("/douyin/message/action/").
		WithQuery("token", tokenA).WithQuery("to_user_id", userIdB).WithQuery("action_type", 1).WithQuery("content", "im top1 and major champion and u?").
		Expect().Status(http.StatusOK).JSON().Object()
	messageResp.Value("status_code").Number().IsEqual(0)

	chatResp := e.GET("/douyin/message/chat/").
		WithQuery("token", tokenA).WithQuery("to_user_id", userIdB).
		Expect().Status(http.StatusOK).JSON().Object()
	chatResp.Value("status_code").Number().IsEqual(0)
	chatResp.Value("message_list").Array().Length().Gt(0)

	chatResp = e.GET("/douyin/message/chat/").
		WithQuery("token", tokenB).WithQuery("to_user_id", userIdA).
		Expect().Status(http.StatusOK).JSON().Object()
	chatResp.Value("status_code").Number().IsEqual(0)
	chatResp.Value("message_list").Array().Length().Gt(0)

}
