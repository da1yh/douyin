package service

import "douyin/dao"

type VideoService interface {
	// AddVideo 把视频添加进数据库
	AddVideo(video dao.Video) error

	// FindPublishedVideosByUserId 通过userId，找到这个人发布的所有视频
	FindPublishedVideosByUserId(userId int64) ([]dao.Video, error)
}
