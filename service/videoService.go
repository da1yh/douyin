package service

import (
	"douyin/dao"
	"time"
)

type VideoService interface {
	// AddVideo 把视频添加进数据库
	AddVideo(video dao.Video) error

	// FindPublishedVideosByUserId 通过userId，找到这个人发布的所有视频
	FindPublishedVideosByUserId(userId int64) ([]dao.Video, error)

	// FindVideosByTimeAndNum 查找发布时间比pubTime早的最近num条视频
	FindVideosByTimeAndNum(pubTime time.Time, num int64) ([]dao.Video, error)

	// FindVideoById 通过视频id，查找video
	FindVideoById(id int64) (dao.Video, error)
}
