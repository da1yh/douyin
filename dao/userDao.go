package dao

import "log"

type User struct {
	Id       int64  `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "user"
}

func FindUserByName(name string) (User, error) {
	user := User{}
	if err := Db.Where("name=?", name).First(&user).Error; err != nil {
		log.Println("error: ", err.Error())
		return user, err
	}
	return user, nil
}

func AddUser(user User) error {
	if err := Db.Create(&user).Error; err != nil {
		log.Println("error: ", err.Error())
		return err
	}
	return nil
}

func FindUserByNameAndPassword(name, password string) (User, error) {
	user := User{}
	if err := Db.Where("name=? AND password=?", name, password).First(&user).Error; err != nil {
		log.Println("error: ", err.Error())
		return user, err
	}
	return user, nil
}

func FindUserById(id int64) (User, error) {
	user := User{}
	if err := Db.Where("id=?", id).First(&user).Error; err != nil {
		log.Println("error: ", err.Error())
		return user, err
	}
	return user, nil
}
