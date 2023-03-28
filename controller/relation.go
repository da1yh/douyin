package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserListResp struct {
	Response
	UserList []User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user not exist",
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
