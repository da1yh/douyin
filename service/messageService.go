package service

import (
	"douyin/dao"
	"time"
)

type MessageService interface {
	// AddMessageByAll 通过所有字段，添加一条消息信息，返回消息的id
	AddMessageByAll(fromUserId, toUserId int64, content string, createTime time.Time) (int64, error)

	// FindMessageIdsByFromUserIdAndToUserId 通过发送人和接受人查找聊天记录（单向）
	FindMessageIdsByFromUserIdAndToUserId(fromUserId, toUserId int64) ([]int64, error)

	// FindMessageById 通过消息id查找消息
	FindMessageById(id int64) (dao.Message, error)
}
