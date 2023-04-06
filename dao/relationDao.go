package dao

import "log"

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
	if err := Db.Where("from_user_id=? AND to_user_id=?", fromUserId, toUserId).First(&relation).Count(&count); err != nil {
		log.Println("error: ", err)
		return count > 0, err.Error
	}
	return count > 0, nil
}
