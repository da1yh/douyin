package dao

import (
	"douyin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

func InitDb() {
	//dsn is data source name
	dsn := config.MysqlUserName + ":" + config.MysqlPassword + "@tcp(" + config.MysqlIp + ":" + config.MysqlPort +
		")/" + config.MysqlDbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln("error: ", err.Error())
	}
	log.Println("connect to mysql database successfully")
}
