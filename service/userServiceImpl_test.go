package service

import (
	"douyin/dao"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	dao.InitDb()
}

//func TestUserServiceImpl_AddUser(t *testing.T) {
//	user := dao.User{
//		Name:     "s1mple",
//		Password: "7355608",
//	}
//	impl := UserServiceImpl{}
//	err := impl.AddUser(user)
//	assert.True(t, err == nil)
//}

func TestUserServiceImpl_FindUserByName(t *testing.T) {
	impl := UserServiceImpl{}
	user, _ := impl.FindUserByName("s1mple")
	assert.True(t, user.Name == "s1mple")
}

func TestUserServiceImpl_FindUserByNameAndPassword(t *testing.T) {
	impl := UserServiceImpl{}
	user, _ := impl.FindUserByNameAndPassword("zywoo", "7355608")
	assert.True(t, user.Name == "zywoo" && user.Password == "7355608")
}

func TestUserServiceImpl_FindUserById(t *testing.T) {
	impl := UserServiceImpl{}
	user, _ := impl.FindUserById(16)
	assert.True(t, user.Name == "niko" && user.Password == "8065537")
}
