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
