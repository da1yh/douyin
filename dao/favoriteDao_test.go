package dao

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit3(t *testing.T) {
	InitDb()
}

func TestFavorite(t *testing.T) {
	err := AddFavoriteByBothId(4, 2)
	assert.Nil(t, err)
	err = AddFavoriteByBothId(14, 2)
	assert.Nil(t, err)
	err = AddFavoriteByBothId(16, 2)
	assert.Nil(t, err)
	num, err := CountFavoritesByToVideoId(2)
	assert.Nil(t, err)
	assert.Equal(t, num, int64(3))
	res, err := CheckFavoriteByBothId(4, 2)
	assert.Nil(t, err)
	assert.True(t, res)
	err = DeleteFavoriteByBothId(4, 2)
	assert.Nil(t, err)
	res, err = CheckFavoriteByBothId(4, 2)
	assert.Nil(t, err)
	assert.False(t, res)
	cnt, err := CountFavoritesByToVideoId(100)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(0)) //count不用加record not found的错误判断

	err = AddFavoriteByBothId(14, 3)
	assert.Nil(t, err)
	ids, err := FindVideoIdsByFromUserId(14)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2)
	ids, err = FindVideoIdsByFromUserId(7)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 0)
}
