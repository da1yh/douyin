package dao

import (
	"douyin/util"
	"log"
	"time"
)

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

func AddMessageByAll(fromUserId, toUserId int64, content string, createTime time.Time) (int64, error) {
	message := Message{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		Content:    content,
		CreateTime: createTime,
	}
	if err := Db.Create(&message).Error; err != nil {
		return util.NotExistId, err
	}
	return message.Id, nil
}

func FindMessageIdsByFromUserIdAndToUserId(fromUserId, toUserId int64) ([]int64, error) {
	ids := make([]int64, 0)
	if err := Db.Model(Message{}).Where("from_user_id=? AND to_user_id=?", fromUserId, toUserId).Pluck("id", &ids).Error; err != nil {
		log.Println(err)
		return ids, err
	}
	return ids, nil
}

func FindMessageById(id int64) (Message, error) {
	message := Message{}
	if err := Db.Model(Message{}).Where("id=?", id).Find(&message).Error; err != nil {
		log.Println(err)
		return message, err
	}
	return message, nil
}
