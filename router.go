package main

import (
	"douyin/controller"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {

	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	//basic api
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	//interact api
	apiRouter.POST("/favorite/action/", controller.Favorite)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.Comment)
	apiRouter.GET("/comment/list/", controller.CommentList)

	//social api
	//apiRouter.POST("/relation/action/", controller.Relation)
	//apiRouter.GET("/relation/follow/list/", controller.RelationFollowList)
	//apiRouter.GET("/relation/follower/list/", controller.RelationFollowerList)
	//apiRouter.GET("/relation/friend/list/", controller.RelationFriendList)
	//apiRouter.GET("/message/chat/", controller.MessageChat)
	//apiRouter.POST("/message/action/", controller.Message)

	//other api

}
