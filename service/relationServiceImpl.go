package service

import (
	"douyin/dao"
	"log"
)

type RelationServiceImpl struct {
}

func (rsi RelationServiceImpl) CountRelationsByFromUserId(id int64) (int64, error) {
	cnt, err := dao.CountRelationsByFromUserId(id)
	if err != nil {
		log.Println("error: ", err.Error())
		log.Println("cannot get follow info")
		return 0, err
	}
	return cnt, nil
}

func (rsi RelationServiceImpl) CountRelationsByToUserId(id int64) (int64, error) {
	cnt, err := dao.CountRelationsByToUserId(id)
	if err != nil {
		log.Println("error: ", err.Error())
		log.Println("cannot get follower info")
		return 0, err
	}
	return cnt, nil
}

func (rsi RelationServiceImpl) CheckRelationByBothId(fromUserId, toUserId int64) (bool, error) {
	flag, err := dao.CheckRelationByBothId(fromUserId, toUserId)
	if err != nil {
		log.Println("error: ", err.Error())
		log.Println("unknown follow relation")
		return false, err
	}
	return flag, nil
}
