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
