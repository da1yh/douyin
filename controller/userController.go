package controller

import (
	"douyin/dao"
	"douyin/middleware/jwt"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
					StatusCode: 2, StatusMsg: "register fail",
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

func UserInfo(c *gin.Context) {
	
	//token := c.Query("token")
	//if user, exist := userLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 0},
	//		User:     user,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "user not exist"},
	//	})
	//}
}
