package dao

import (
	"gorm.io/gorm"
	"log"
)

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
	//var count int64
	if err := Db.Where("from_user_id=? AND to_video_id=?", fromUserId, toVideoId).First(&favorite).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		log.Println("err: ", err)
		return false, err
	}
	return true, nil
}

func AddFavoriteByBothId(fromUserId, toVideoId int64) error {
	favorite := Favorite{
		FromUserId: fromUserId,
		ToVideoId:  toVideoId,
	}
	if err := Db.Create(&favorite).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteFavoriteByBothId(fromUserId, toVideoId int64) error {
	//favorite := Favorite{
	//	FromUserId: fromUserId,
	//	ToVideoId:  toVideoId,
	//}
	if err := Db.Where("from_user_id=? AND to_video_id=?", fromUserId, toVideoId).Delete(Favorite{}).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}
