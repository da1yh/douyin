package test

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

var serverAddr = "http://localhost:8080"

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
		BaseURL:  serverAddr,
		Reporter: httpexpect.Reporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
