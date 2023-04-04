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
