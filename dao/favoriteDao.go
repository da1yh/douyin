package dao

import "log"

type Favorite struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToVideoId  int64 `gorm:"column:to_video_id"`
}

func (Favorite) TableName() string {
	return "favorite"
}

func CountFavoritesByToVideoId(toVideoId int64) (int64, error) {
	favorite := Favorite{}
	var count int64
	if err := Db.Where("to_video_id=?", toVideoId).Find(&favorite).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func CheckFavoriteByBothId(fromUserId, toVideoId int64) (bool, error) {
	favorite := Favorite{}
	var count int64
	if err := Db.Where("from_user_id=? AND to_video_id=?", fromUserId, toVideoId).First(&favorite).Count(&count).Error; err != nil {
		log.Println("err: ", err)
		return false, err
	}
	return count > 0, nil
}
