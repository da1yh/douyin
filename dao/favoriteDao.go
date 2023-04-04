package dao

type Favorite struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToVideoId  int64 `gorm:"column:to_video_id"`
}

func (Favorite) TableName() string {
	return "favorite"
}
