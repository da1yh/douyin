package dao

import "time"

type Message struct {
	Id         int64     `gorm:"column:id"`
	FromUserId int64     `gorm:"column:from_user_id"`
	ToUserId   int64     `gorm:"column:to_user_id"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (Message) TableName() string {
	return "message"
}
