package dao

import (
	"gorm.io/gorm"
	"log"
)

type Relation struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToUserId   int64 `gorm:"column:to_user_id"`
}

func (Relation) TableName() string {
	return "relation"
}

func CountRelationsByFromUserId(id int64) (int64, error) {
	relation := Relation{}
	var count int64
	err := Db.Where("from_user_id=?", id).Find(&relation).Count(&count)
	if err != nil {
		log.Println("error: ", err)
		return count, err.Error
	}
	return count, nil
}

func CountRelationsByToUserId(id int64) (int64, error) {
	relation := Relation{}
	var count int64
	if err := Db.Where("to_user_id=?", id).Find(&relation).Count(&count); err != nil {
		log.Println("error: ", err)
		return count, err.Error
	}
	return count, nil
}

func CheckRelationByBothId(fromUserId, toUserId int64) (bool, error) {
	relation := Relation{}
	var count int64
	if err := Db.Where("from_user_id=? AND to_user_id=?", fromUserId, toUserId).First(&relation).Count(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

// FindRelationFriendIdsByFromUserId 通过fromUserId找到它的朋友的id列表
func FindRelationFriendIdsByFromUserId(fromUserId int64) ([]int64, error) {
	friendIds := make([]int64, 0)
	if err := Db.Table("relation r1").
		Select("r1.to_user_id").
		Joins("JOIN relation r2 ON r1.from_user_id = r2.to_user_id AND r1.to_user_id = r2.from_user_id").
		Where("r1.from_user_id=? AND r1.from_user_id < r1.to_user_id", fromUserId).
		Find(&friendIds).Error; err != nil {
		return friendIds, err
	}
	return friendIds, nil
}

func AddRelationByBothId(fromUserID, toUserId int64) error {
	relation := Relation{
		FromUserId: fromUserID,
		ToUserId:   toUserId,
	}
	if err := Db.Create(&relation).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteRelationByBothId(fromUserId, toUserId int64) error {
	if err := Db.Model(&Relation{}).Where("from_user_id=? AND to_user_id=?", fromUserId, toUserId).Delete(&Relation{}).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// FindRelationToUserIdsByFromUserId 找到某个人关注的用户id列表
func FindRelationToUserIdsByFromUserId(fromUserId int64) ([]int64, error) {
	ids := make([]int64, 0)
	if err := Db.Model(&Relation{}).Where("from_user_id=?", fromUserId).Pluck("to_user_id", &ids).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// FindRelationFromUserIdsByToUserId 找到某个人的粉丝id列表
func FindRelationFromUserIdsByToUserId(toUserId int64) ([]int64, error) {
	ids := make([]int64, 0)
	if err := Db.Model(&Relation{}).Where("to_user_id=?", toUserId).Pluck("from_user_id", &ids).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}
