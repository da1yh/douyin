package dao

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 4->1 "good video" X
// 4->2 "good one"
// 14->2 "hello you" X
// 4->2 "haha"
// 14->1 "yes"
func TestCommentDao(t *testing.T) {
	InitDb()
	id1, err := AddCommentByAll(4, 1, "good video", time.Now())
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	id2, err := AddCommentByAll(4, 2, "good one", time.Now())
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	id3, err := AddCommentByAll(14, 2, "hello you", time.Now())
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	id4, err := AddCommentByAll(4, 2, "haha", time.Now())
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	id5, err := AddCommentByAll(14, 1, "yes", time.Now())
	assert.Nil(t, err)

	ids, err := FindCommentIdsByToVideoId(2)
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 3)
	a, b, c := 0, 0, 0
	for _, id := range ids {
		if id == id2 {
			a++
		}
		if id == id3 {
			b++
		}
		if id == id4 {
			c++
		}
	}
	assert.Equal(t, a, 1)
	assert.Equal(t, b, 1)
	assert.Equal(t, c, 1)

	id, err := FindCommentToVideoIdById(id1)
	assert.Nil(t, err)
	assert.Equal(t, id, int64(1))
	id, err = FindCommentToVideoIdById(id2)
	assert.Nil(t, err)
	assert.Equal(t, id, int64(2))
	id, err = FindCommentToVideoIdById(id3)
	assert.Nil(t, err)
	assert.Equal(t, id, int64(2))
	id, err = FindCommentToVideoIdById(id4)
	assert.Nil(t, err)
	assert.Equal(t, id, int64(2))
	id, err = FindCommentToVideoIdById(id5)
	assert.Nil(t, err)
	assert.Equal(t, id, int64(1))

	err = DeleteCommentById(id1)
	assert.Nil(t, err)
	err = DeleteCommentById(id3)
	assert.Nil(t, err)
	//然后查数据库对不对

}
