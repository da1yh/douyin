package service

import (
	"douyin/dao"
	"log"
)

type CommentServiceImpl struct {
}

func (csi CommentServiceImpl) CountCommentsByToVideoId(toVideoId int64) (int64, error) {
	cnt, err := dao.CountCommentsByToVideoId(toVideoId)
	if err != nil {
		log.Println("err: ", err)
		log.Println("get the number of video's comments failed")
	}
	return cnt, nil
}
