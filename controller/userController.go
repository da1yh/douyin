package controller

import (
	"douyin/dao"
	"douyin/middleware/jwt"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// map that key is token, value is UserInfo, for checking if user already exist when register and login
var userLoginInfo = map[string]User{
	"dyh_12345": {
		Id:            1,
		Name:          "dyh",
		FollowCount:   15,
		FollowerCount: 16,
		IsFollow:      false,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user,omitempty"`
}

var userIdSequence int64 = 1

// UserRegister
// 获得用户名和密码，调用userService的方法查看用户是否存在，根据不同的情况进行返回/*
func UserRegister(c *gin.Context) {
	username, password := c.Query("username"), c.Query("password")
	usi := service.UserServiceImpl{}
	user, _ := usi.FindUserByName(username)
	if user.Name == username {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "user already exist"},
		})
	} else {
		newUser := dao.User{
			Name:     username,
			Password: password,
		}
		err := usi.AddUser(newUser)
		if err != nil {
			log.Println("error: ", err.Error())
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 2, StatusMsg: "register failed",
				},
			})
		}
		user, err = usi.FindUserByName(username)
		if err != nil {
			log.Println("error: ", err.Error())
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0, StatusMsg: "register successfully",
			},
			UserId: user.Id,
			Token:  jwt.GenerateToken(user.Id, user.Name, user.Password),
		})
	}
}

// UserLogin 获得用户名和密码，查看是否存在相应的值，根据不同的情况进行返回
func UserLogin(c *gin.Context) {
	username, password := c.Query("username"), c.Query("password")
	usi := service.UserServiceImpl{}
	user, _ := usi.FindUserByNameAndPassword(username, password)
	if user.Name != username {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1, StatusMsg: "either username or password is wrong",
			},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0, StatusMsg: "login successfully",
			},
			UserId: user.Id,
			Token:  jwt.GenerateToken(user.Id, user.Name, user.Password),
		})
	}
}

// UserInfo 鉴权通过后，通过id，在数据库查用户是否存在，如果存在，在relation表中查follow相关信息，组成返回结构体
func UserInfo(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	usi := service.UserServiceImpl{}
	user, err := usi.FindUserById(userId)
	if err != nil || user.Id != userId {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "user not found",
		})
	} else {
		rsi := service.RelationServiceImpl{}
		var followCount, followerCount int64
		var isFollow bool
		followCount, _ = rsi.CountRelationsByFromUserId(userId)
		followerCount, _ = rsi.CountRelationsByToUserId(userId)
		curId := c.GetInt64("id")
		isFollow, _ = rsi.CheckRelationByBothId(curId, userId)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 0, StatusMsg: "get user info successfully",
			},
			User: User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   followCount,
				FollowerCount: followerCount,
				IsFollow:      isFollow,
			},
		})
	}
}
