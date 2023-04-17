package service

import (
	"douyin/dao"
	"log"
)

type FavoriteServiceImpl struct {
}

func (fsi FavoriteServiceImpl) CountFavoritesByToVideoId(toVideoId int64) (int64, error) {
	cnt, err := dao.CountFavoritesByToVideoId(toVideoId)
	if err != nil {
		log.Println("err: ", err.Error())
		return 0, err
	}
	return cnt, nil
}

func (fsi FavoriteServiceImpl) CheckFavoriteByBothId(fromUserId, toVideoId int64) (bool, error) {
	flag, err := dao.CheckFavoriteByBothId(fromUserId, toVideoId)
	if err != nil {
		log.Println("err: ", err)
		return false, err
	}
	return flag, nil
}
