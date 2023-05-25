package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResp struct {
	Response
	UserList []User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	fromUserId, toUserIdStr, actionType := c.GetInt64("id"), c.Query("to_user_id"), c.Query("action_type")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "failed to get to_user_id",
		})
	}
	rsi := service.RelationServiceImpl{}
	if actionType == "1" {
		err = rsi.AddRelationByBothId(fromUserId, toUserId)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1, StatusMsg: "failed to follow this user",
			})
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0, StatusMsg: "follow this user successfully",
		})
	} else {
		err = rsi.DeleteRelationByBothId(fromUserId, toUserId)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 2, StatusMsg: "failed to unfollow this user",
			})
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0, StatusMsg: "unfollow this user successfully",
		})
	}
}

func RelationFollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResp{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

func RelationFollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResp{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

func RelationFriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResp{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
