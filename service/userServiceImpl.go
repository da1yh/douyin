package service

import (
	"douyin/dao"
	"log"
)

type UserServiceImpl struct {
}

func (usi UserServiceImpl) FindUserByName(name string) (dao.User, error) {
	user, err := dao.FindUserByName(name)
	if err != nil {
		log.Println("error: ", err.Error())
		log.Println("user not found")
		return user, err
	}
	return user, nil
}

func (usi UserServiceImpl) AddUser(user dao.User) error {
	if err := dao.AddUser(user); err != nil {
		return err
	}
	return nil
}

func (usi UserServiceImpl) FindUserByNameAndPassword(name, password string) (dao.User, error) {
	user, err := dao.FindUserByNameAndPassword(name, password)
	if err != nil {
		log.Println("error: ", err.Error())
		log.Println("user not found")
		return user, err
	}
	return user, nil
}
