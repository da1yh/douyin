package dao

import "time"

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
