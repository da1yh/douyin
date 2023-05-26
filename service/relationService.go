package service

type RelationService interface {
	// CountRelationsByFromUserId 通过userid，查找他关注了多少人
	CountRelationsByFromUserId(fromUserId int64) (int64, error)

	// CountRelationsByToUserId 通过userid，查找他的粉丝数
	CountRelationsByToUserId(toUserId int64) (int64, error)

	// CheckRelationByBothId 通过双方的id，查找是否存在关注关系
	CheckRelationByBothId(fromUserId, toUserId int64) (bool, error)

	// AddRelationByBothId 通过双方的id，添加关注关系
	AddRelationByBothId(fromUserId, toUserId int64) error

	// DeleteRelationByBothId 通过双方的id，删除关注关系
	DeleteRelationByBothId(fromUserId, toUserId int64) error

	// FindRelationToUserIdsByFromUserId 找到该用户关注的人的id列表
	FindRelationToUserIdsByFromUserId(fromUserId int64) ([]int64, error)

	// FindRelationFromUserIdsByToUserId 找到该用户的粉丝id列表
	FindRelationFromUserIdsByToUserId(toUserId int64) ([]int64, error)
}
