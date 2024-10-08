package main

import (
	"douyin/controller"
	"douyin/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {

	//r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	//basic api
	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.UserRegister)
	apiRouter.POST("/user/login/", controller.UserLogin)

	apiRouter.GET("/feed/", jwt.AuthVisitor(), controller.Feed)
	apiRouter.POST("/publish/action/", jwt.AuthPost(), controller.PublishAction)
	apiRouter.GET("/publish/list/", jwt.Auth(), controller.PublishList)

	//interact api
	apiRouter.POST("/favorite/action/", jwt.AuthPost(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", jwt.Auth(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", jwt.AuthPost(), controller.CommentAction)
	apiRouter.GET("/comment/list/", jwt.Auth(), controller.CommentList)

	//social api
	apiRouter.POST("/relation/action/", jwt.AuthPost(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", jwt.Auth(), controller.RelationFollowList)
	apiRouter.GET("/relation/follower/list/", jwt.Auth(), controller.RelationFollowerList)
	apiRouter.GET("/relation/friend/list/", jwt.Auth(), controller.RelationFriendList)
	apiRouter.GET("/message/chat/", jwt.Auth(), controller.MessageChat)
	apiRouter.POST("/message/action/", jwt.AuthPost(), controller.MessageAction)

}
