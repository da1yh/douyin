package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
)

var userLoginInfo = map[string]User{
	"dyh_12345": {
		Id:            1,
		Name:          "dyh",
		FollowCount:   15,
		FollowerCount: 16,
		IsFollow:      false,
	},
}

type UserLoginResp struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

var userIdSequence int64 = 1

func Register(c *gin.Context) {
	//get username and password, check if user already exists, if not, respond new user info
	username, password := c.Query("username"), c.Query("password")
	token := username + "_" + password
	if _, exist := userLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResp{
			Response: Response{StatusCode: 1, StatusMsg: "user already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Id:   userIdSequence,
			Name: username,
		}
		userLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResp{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {

}

func UserInfo(c *gin.Context) {

}
