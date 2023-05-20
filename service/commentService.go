package service

import (
	"douyin/dao"
	"time"
)

type CommentService interface {
	// CountCommentsByToVideoId 通过videoId查询该视频的评论数
	CountCommentsByToVideoId(toVideoId int64) (int64, error)

	// AddCommentByAll 通过给定所有字段的信息，增加一条评论信息，返回commentId
	AddCommentByAll(fromUserId, toVideoId int64, content string, createDate time.Time) (int64, error)

	// DeleteCommentById 通过评论id删除评论，返回该评论
	DeleteCommentById(id int64) (dao.Comment, error)

	// FindCommentIdsByToVideoId 通过toVideoId返回该视频的评论id列表
	FindCommentIdsByToVideoId(toVideoId int64) ([]int64, error)

	// FindCommentById 通过id找到comment
	FindCommentById(id int64) (dao.Comment, error)
}
