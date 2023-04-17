package service

type CommentService interface {
	// 通过videoid查询该视频的评论数
	CountCommentsByToVideoId(toVideoId int64) (int64, error)
}
