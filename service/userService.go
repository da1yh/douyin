package service

import "douyin/dao"

type UserService interface {
	// FindUserByName 根据用户名找到用户，err表示找不到
	FindUserByName(name string) (dao.User, error)

	// AddUser 把用户添加进数据库
	AddUser(user dao.User) error

	// FindUserByNameAndPassword 根据用户名和密码找到用户
	FindUserByNameAndPassword(name, password string) (dao.User, error)

	// FindUserById 根据id查找用户
	FindUserById(id int64) (dao.User, error)
}
