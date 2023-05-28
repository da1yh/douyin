package controller

import (
	"douyin/config"
	"douyin/service"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title,omitempty"`
}

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	FromUserId int64  `json:"from_user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"from_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

func getTestUser(user string, e *httpexpect.Expect) (int, string) {
	//register using "user" as username and password, if registered, login
	registerResp := e.POST("/douyin/user/register/").
		WithQuery("username", user).WithQuery("password", user).
		WithFormField("username", user).WithFormField("password", user).
		Expect().Status(http.StatusOK).JSON().Object()
	userId, token := 0, ""
	token = registerResp.Value("token").String().Raw()
	if len(token) == 0 { // user already exist, just login
		loginResp := e.POST("/douyin/user/login/").
			WithQuery("username", user).WithQuery("password", user).
			WithFormField("username", user).WithFormField("password", user).
			Expect().Status(http.StatusOK).JSON().Object()

		loginResp.Value("token").String().Length().Gt(0)
		token = loginResp.Value("token").String().Raw()
		userId = int(loginResp.Value("user_id").Number().Raw())
	} else {
		userId = int(registerResp.Value("user_id").Number().Raw())
	}
	return userId, token
}

func newExpect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client:   http.DefaultClient,
		BaseURL:  config.ServerAddr,
		Reporter: httpexpect.Reporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

// GetUserRespByBothId 通过当前用户的id和待查询用户的id，获得UserResponse，游客id为myId=-1
func GetUserRespByBothId(myId, yourId int64) (User, error) {
	rsi := service.RelationServiceImpl{}
	usi := service.UserServiceImpl{}
	var followCount, followerCount int64
	var isFollow bool
	followCount, _ = rsi.CountRelationsByFromUserId(yourId)
	followerCount, _ = rsi.CountRelationsByToUserId(yourId)
	if myId == -1 {
		isFollow = false
	} else {
		isFollow, _ = rsi.CheckRelationByBothId(myId, yourId)
	}
	usr, _ := usi.FindUserById(yourId)
	user := User{
		Id:            yourId,
		Name:          usr.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}
	return user, nil
}

// GetTokenByLogin make sure this user already registered
func GetTokenByLogin(username, password string, e *httpexpect.Expect) string {
	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", username).WithQuery("password", password).
		Expect().Status(http.StatusOK).JSON().Object()
	loginResp.Value("status_code").Number().IsEqual(0)
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().Length().Gt(0)
	token := loginResp.Value("token").String().Raw()
	return token
}
