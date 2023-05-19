package dao

import (
	"douyin/util"
	"log"
	"time"
)

type Comment struct {
	Id         int64     `gorm:"column:id"`
	FromUserId int64     `gorm:"column:from_user_id"`
	ToVideoId  int64     `gorm:"column:to_video_id"`
	Content    string    `gorm:"column:content"`
	CreateDate time.Time `gorm:"column:create_date"`
}

func (Comment) TableName() string {
	return "comment"
}

func CountCommentsByToVideoId(toVideoId int64) (int64, error) {
	comment := Comment{}
	var count int64
	if err := Db.Where("to_video_id=?", toVideoId).Find(&comment).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func AddCommentByAll(fromUserId, toVideoId int64, content string, createDate time.Time) (int64, error) {
	comment := Comment{
		FromUserId: fromUserId,
		ToVideoId:  toVideoId,
		Content:    content,
		CreateDate: createDate,
	}
	if err := Db.Create(&comment).Error; err != nil {
		return util.NotExistId, err
	}
	return comment.Id, nil
}

func DeleteCommentByAll(fromUserId, toVideoId int64, content string, createDate time.Time) error {
	if err := Db.Where("from_user_id=? AND to_video_id=? AND content=? AND create_date=?", fromUserId, toVideoId, content, createDate).Delete(Comment{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCommentById(id int64) error {
	if err := Db.Model(Comment{}).Where("id=?", id).Delete(Comment{}).Error; err != nil {
		return err
	}
	return nil
}

// FindCommentIdsByToVideoId 根据toVideoId找到所有的评论，返回id字段的值
func FindCommentIdsByToVideoId(toVideoId int64) ([]int64, error) {
	ids := make([]int64, 0)
	if err := Db.Model(Comment{}).Where("to_video_id=?", toVideoId).Pluck("id", &ids).Error; err != nil {
		log.Println(err)
		return ids, err
	}
	return ids, nil
}

func FindCommentToVideoIdById(id int64) (int64, error) {
	var toVideoId int64 //此处pluck能用非slice？
	if err := Db.Model(Comment{}).Where("id=?", id).Pluck("to_video_id", &toVideoId).Error; err != nil {
		log.Println(err)
		return toVideoId, err
	}
	return toVideoId, nil
}
