package service

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// user 4 14 16 17
// 4 -> 14
// 4 -> 16
// 14 -> 16
// 16 -> 4 X
// 17 -> 4
// 4 -> 17 X
// 16 -> 17
// [17 -> 16]
// add check(14 -> 16 16 -> 4 17 -> 16)  find ( 16's follow 17's follower 4's friend)
// delete add check (14 -> 16 16 -> 4 17 -> 16) find (16's follow 17's follower 4's friend)
func TestRelationServiceImpl(t *testing.T) {
	dao.InitDb()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	forever := make(chan struct{})
	rabbitmq.InitRelationMQ()

	time.Sleep(3 * time.Second)

	rsi := RelationServiceImpl{}
	err := rsi.AddRelationByBothId(4, 14)
	assert.Nil(t, err)
	err = rsi.AddRelationByBothId(4, 16)
	assert.Nil(t, err)
	err = rsi.AddRelationByBothId(14, 16)
	assert.Nil(t, err)
	err = rsi.AddRelationByBothId(16, 4)
	assert.Nil(t, err)
	err = rsi.AddRelationByBothId(17, 4)
	assert.Nil(t, err)
	err = rsi.AddRelationByBothId(4, 17)
	assert.Nil(t, err)
	err = rsi.AddRelationByBothId(16, 17)
	assert.Nil(t, err)

	res, err := rsi.CheckRelationByBothId(14, 16)
	assert.Nil(t, err)
	assert.True(t, res)
	res, err = rsi.CheckRelationByBothId(16, 4)
	assert.Nil(t, err)
	assert.True(t, res)
	res, err = rsi.CheckRelationByBothId(17, 16)
	assert.Nil(t, err)
	assert.False(t, res)

	ids, err := rsi.FindRelationToUserIdsByFromUserId(16)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2)
	assert.True(t, (ids[0] == int64(4) && ids[1] == int64(17)) || (ids[0] == int64(17) && ids[1] == int64(4)))
	cnt, err := rsi.CountRelationsByFromUserId(16)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(2))

	ids, err = rsi.FindRelationFromUserIdsByToUserId(17)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2)
	assert.True(t, (ids[0] == int64(4) && ids[1] == int64(16)) || (ids[0] == int64(16) && ids[1] == int64(4)))
	cnt, err = rsi.CountRelationsByToUserId(17)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(2))

	ids, err = rsi.FindRelationFriendIdsByFromUserId(4)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2)
	assert.True(t, (ids[0] == int64(16) && ids[1] == int64(17)) || (ids[0] == int64(17) && ids[1] == int64(16)))

	err = rsi.DeleteRelationByBothId(16, 4)
	assert.Nil(t, err)
	err = rsi.DeleteRelationByBothId(4, 17)
	assert.Nil(t, err)
	err = rsi.AddRelationByBothId(17, 16)
	assert.Nil(t, err)

	res, err = rsi.CheckRelationByBothId(14, 16)
	assert.Nil(t, err)
	assert.True(t, res)
	res, err = rsi.CheckRelationByBothId(16, 4)
	assert.Nil(t, err)
	assert.False(t, res)
	res, err = rsi.CheckRelationByBothId(17, 16)
	assert.Nil(t, err)
	assert.True(t, res)

	ids, err = rsi.FindRelationToUserIdsByFromUserId(16)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 1)
	assert.Equal(t, ids[0], int64(17))
	cnt, err = rsi.CountRelationsByFromUserId(16)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(1))

	ids, err = rsi.FindRelationFromUserIdsByToUserId(17)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 1)
	assert.Equal(t, ids[0], int64(16))
	cnt, err = rsi.CountRelationsByToUserId(17)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(1))

	ids, err = rsi.FindRelationFriendIdsByFromUserId(4)
	assert.Nil(t, err)
	assert.Empty(t, ids)

	// check redis correct

	<-forever
}
