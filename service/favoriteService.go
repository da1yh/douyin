package service

type FavoriteService interface {
	// CountFavoritesByToVideoId 通过videoid，查找被多少人点赞
	CountFavoritesByToVideoId(toVideoId int64) (int64, error)

	// CheckFavoriteByBothId 查询某个人是否点赞了该视频
	CheckFavoriteByBothId(fromUserId, toVideoId int64) (bool, error)
}
