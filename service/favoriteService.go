package service

type FavoriteService interface {
	// CountFavoritesByToVideoId 通过videoId，查找被多少人点赞
	CountFavoritesByToVideoId(toVideoId int64) (int64, error)

	// CheckFavoriteByBothId 查询某个人是否点赞了该视频
	CheckFavoriteByBothId(fromUserId, toVideoId int64) (bool, error)

	// AddFavoriteByBothId 通过id增加某个人点赞某个视频的记录
	AddFavoriteByBothId(fromUserId, toVideoId int64) error

	// DeleteFavoriteByBothId 通过id删除某个人点赞某个视频的记录
	DeleteFavoriteByBothId(fromUserId, toVideoId int64) error

	// FindFavoriteVideoIdsByFromUserId 通过fromUserId，查找该用户点赞的视频的id列表
	FindFavoriteVideoIdsByFromUserId(fromUserId int64) ([]int64, error)
}
