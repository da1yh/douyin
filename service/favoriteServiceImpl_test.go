package service

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFavoriteServiceImpl(t *testing.T) {
	dao.InitDb()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFavoriteMQ()
	fsi := FavoriteServiceImpl{}
	err := fsi.AddFavoriteByBothId(4, 1)
	assert.Nil(t, err)
	err = fsi.AddFavoriteByBothId(4, 2)
	assert.Nil(t, err)
	err = fsi.AddFavoriteByBothId(4, 3)
	assert.Nil(t, err)
	err = fsi.AddFavoriteByBothId(14, 3)
	assert.Nil(t, err)
	err = fsi.AddFavoriteByBothId(16, 3)
	assert.Nil(t, err)

	res, err := fsi.FindFavoriteVideoIdsByFromUserId(4)
	assert.Nil(t, err)
	assert.Equal(t, len(res), 3)
	one, two, three := 0, 0, 0
	for _, v := range res {
		if v == 1 {
			one++
		} else if v == 2 {
			two++
		} else if v == 3 {
			three++
		}
	}
	assert.Equal(t, one, 1)
	assert.Equal(t, two, 1)
	assert.Equal(t, three, 1)

	cnt, err := fsi.CountFavoritesByToVideoId(3)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(3))

	exist, err := fsi.CheckFavoriteByBothId(4, 3)
	assert.Nil(t, err)
	assert.True(t, exist)
	exist, err = fsi.CheckFavoriteByBothId(14, 2)
	assert.Nil(t, err)
	assert.False(t, exist)

	err = fsi.DeleteFavoriteByBothId(4, 3)
	assert.Nil(t, err)
	cnt, err = fsi.CountFavoritesByToVideoId(3)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(2))
	exist, err = fsi.CheckFavoriteByBothId(4, 3)
	assert.Nil(t, err)
	assert.False(t, exist)
}

func TestFavoriteServiceImpl_FindFavoriteVideoIdsByFromUserId(t *testing.T) {
	dao.InitDb()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFavoriteMQ()
	fsi := FavoriteServiceImpl{}
	ids, err := fsi.FindFavoriteVideoIdsByFromUserId(14)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2)
	assert.True(t, (ids[0] == int64(1) && ids[1] == int64(2)) || (ids[0] == int64(2) && ids[1] == int64(1)))
}
