package service

import (
	"douyin/dao"
	"time"
)

type VideoServiceImpl struct {
}

func (vsi VideoServiceImpl) AddVideo(video dao.Video) error {
	if err := dao.AddVideo(video); err != nil {
		return err
	}
	return nil
}

func (vsi VideoServiceImpl) FindPublishedVideosByUserId(userId int64) ([]dao.Video, error) {
	videos := make([]dao.Video, 0)
	videos, err := dao.FindPublishedVideosByUserId(userId)
	if err != nil {
		return videos, err
	}
	return videos, nil
}

func (vsi VideoServiceImpl) FindVideosByTimeAndNum(pubTime time.Time, num int) ([]dao.Video, error) {
	videos := make([]dao.Video, 0, num)
	videos, err := dao.FindVideosByTimeAndNum(pubTime, num)
	if err != nil {
		return videos, err
	}
	return videos, nil
}

func (vsi VideoServiceImpl) FindVideoById(id int64) (dao.Video, error) {
	video, err := dao.FindVideoById(id)
	if err != nil {
		return video, err
	}
	return video, nil
}
