package dao

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit1(t *testing.T) {
	InitDb()
	friendIds, err := FindRelationFriendIdsByFromUserId(4)
	assert.Nil(t, err)
	assert.Equal(t, len(friendIds), 2)
	a, b := 0, 0
	for _, friendId := range friendIds {
		if friendId == 14 {
			a++
		} else if friendId == 16 {
			b++
		}
	}
	assert.Equal(t, a, 1)
	assert.Equal(t, b, 1)
}
