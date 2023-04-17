package dao

import (
	"log"
	"time"
)

type Video struct {
	Id          int64     `gorm:"column:id"`
	UserId      int64     `gorm:"column:user_id"`
	PublishTime time.Time `gorm:"column:publish_time"`
	PlayUrl     string    `gorm:"column:play_url"`
	CoverUrl    string    `gorm:"column:cover_url"`
	Title       string    `gorm:"column:title"`
}

func (Video) TableName() string {
	return "video"
}

func AddVideo(video Video) error {
	if err := Db.Create(&video).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}
