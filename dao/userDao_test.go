package dao

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	InitDb()
}

//func TestAddUser(t *testing.T) {
//	user := User{
//		Name:     "zywoo",
//		Password: "7355608",
//	}
//	err := AddUser(user)
//	if err != nil {
//		fmt.Println("error: ", err.Error())
//	}
//}

func TestFindUserByName(t *testing.T) {
	user, err := FindUserByName("s1mple")
	assert.True(t, user.Name == "s1mple" && user.Password == "7355608")
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
}

func TestFindUserByNameAndPassword(t *testing.T) {
	user, _ := FindUserByNameAndPassword("zywoo", "7355608")
	assert.True(t, user.Name == "zywoo" && user.Password == "7355608")
}

func TestFindUserById(t *testing.T) {
	user, _ := FindUserById(16)
	assert.True(t, user.Name == "niko" && user.Password == "8065537")
}
