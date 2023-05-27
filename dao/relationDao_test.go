package dao

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
func TestInit1(t *testing.T) {
	InitDb()
	err := AddRelationByBothId(4, 14)
	assert.Nil(t, err)
	err = AddRelationByBothId(4, 16)
	assert.Nil(t, err)
	err = AddRelationByBothId(14, 16)
	assert.Nil(t, err)
	err = AddRelationByBothId(16, 4)
	assert.Nil(t, err)
	err = AddRelationByBothId(17, 4)
	assert.Nil(t, err)
	err = AddRelationByBothId(4, 17)
	assert.Nil(t, err)
	err = AddRelationByBothId(16, 17)
	assert.Nil(t, err)

	res, err := CheckRelationByBothId(14, 16)
	assert.Nil(t, err)
	assert.True(t, res)
	res, err = CheckRelationByBothId(16, 4)
	assert.Nil(t, err)
	assert.True(t, res)
	res, err = CheckRelationByBothId(17, 16)
	assert.Nil(t, err)
	assert.False(t, res)

	ids, err := FindRelationToUserIdsByFromUserId(16)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2)
	assert.True(t, (ids[0] == int64(4) && ids[1] == int64(17)) || (ids[0] == int64(17) && ids[1] == int64(4)))

	ids, err = FindRelationFromUserIdsByToUserId(17)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2)
	assert.True(t, (ids[0] == int64(4) && ids[1] == int64(16)) || (ids[0] == int64(16) && ids[1] == int64(4)))

	ids, err = FindRelationFriendIdsByFromUserId(4)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2)
	assert.True(t, (ids[0] == int64(16) && ids[1] == int64(17)) || (ids[0] == int64(17) && ids[1] == int64(16)))

	err = DeleteRelationByBothId(16, 4)
	assert.Nil(t, err)
	err = DeleteRelationByBothId(4, 17)
	assert.Nil(t, err)
	err = AddRelationByBothId(17, 16)
	assert.Nil(t, err)

	res, err = CheckRelationByBothId(14, 16)
	assert.Nil(t, err)
	assert.True(t, res)
	res, err = CheckRelationByBothId(16, 4)
	assert.Nil(t, err)
	assert.False(t, res)
	res, err = CheckRelationByBothId(17, 16)
	assert.Nil(t, err)
	assert.True(t, res)

	ids, err = FindRelationToUserIdsByFromUserId(16)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 1)
	assert.Equal(t, ids[0], int64(17))

	ids, err = FindRelationFromUserIdsByToUserId(17)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 1)
	assert.Equal(t, ids[0], int64(16))

	ids, err = FindRelationFriendIdsByFromUserId(4)
	assert.Nil(t, err)
	assert.Empty(t, ids)

}
