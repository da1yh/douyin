package service

import "douyin/dao"

type VideoService interface {
	// AddVideo 把视频添加进数据库
	AddVideo(video dao.Video) error
}
