package service

import (
	"douyin/dao"
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
