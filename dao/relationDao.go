package dao

type Relation struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToUserId   int64 `gorm:"column:to_user_id"`
}

func (Relation) TableName() string {
	return "relation"
}
