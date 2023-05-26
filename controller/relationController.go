package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"log"
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
	myId := c.GetInt64("id")
	fromUserIdStr := c.Query("user_id")
	fromUserId, err := strconv.ParseInt(fromUserIdStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "failed to parse user id",
		})
	}
	rsi := service.RelationServiceImpl{}
	toUserIds, err := rsi.FindRelationToUserIdsByFromUserId(fromUserId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 2, StatusMsg: "failed to find idols in database",
		})
	}
	userList := make([]User, 0)
	for _, toUserId := range toUserIds {
		userResp, err := GetUserRespByBothId(myId, toUserId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 3, StatusMsg: "failed to find user info",
			})
		}
		userList = append(userList, userResp)
	}
	c.JSON(http.StatusOK, UserListResp{
		Response: Response{
			StatusCode: 0, StatusMsg: "get follow list successfully",
		},
		UserList: userList,
	})
}

func RelationFollowerList(c *gin.Context) {
	myId := c.GetInt64("id")
	toUserIdStr := c.Query("user_id")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "failed to parse user id",
		})
	}
	rsi := service.RelationServiceImpl{}
	fromUserIds, err := rsi.FindRelationFromUserIdsByToUserId(toUserId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 2, StatusMsg: "failed to find fans in database",
		})
	}
	userList := make([]User, 0)
	for _, fromUserId := range fromUserIds {
		userResp, err := GetUserRespByBothId(myId, fromUserId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 3, StatusMsg: "failed to find user info",
			})
		}
		userList = append(userList, userResp)
	}
	c.JSON(http.StatusOK, UserListResp{
		Response: Response{
			StatusCode: 0, StatusMsg: "get follower list successfully",
		},
		UserList: userList,
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
