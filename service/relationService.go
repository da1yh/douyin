package service

type RelationService interface {
	// CountRelationsByFromUserId 通过userid，查找他关注了多少人
	CountRelationsByFromUserId(id int64) (int64, error)

	// CountRelationsByToUserId 通过userid，查找他的粉丝数
	CountRelationsByToUserId(id int64) (int64, error)

	// CheckRelationByBothId 通过双方的id，查找是否存在关注关系
	CheckRelationByBothId(fromUserId, toUserId int64) (bool, error)
}
