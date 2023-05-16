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

func FindPublishedVideosByUserId(userId int64) ([]Video, error) {
	videos := make([]Video, 0)
	if err := Db.Where("user_id=?", userId).Find(&videos).Error; err != nil {
		log.Println(err)
		return videos, err
	}
	return videos, nil
}

func FindVideosByTimeAndNum(pubTime time.Time, num int) ([]Video, error) {
	videos := make([]Video, 0, num)
	if err := Db.Where("publish_time<?", pubTime).Order("publish_time desc").Limit(num).Find(&videos).Error; err != nil {
		log.Println(err)
		return videos, err
	}
	return videos, nil
}

func FindVideoById(id int64) (Video, error) {
	video := Video{}
	if err := Db.Where("id=?", id).First(&video).Error; err != nil {
		return video, err
	}
	return video, nil
}
