package service

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 4->1 "good video" X
// 4->2 "good one"
// 14->2 "hello you" X
// 4->2 "haha"
// 14->1 "yes"
func TestCommentServiceImpl(t *testing.T) {
	dao.InitDb()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	forever := make(chan struct{})
	rabbitmq.InitCommentMQ()

	time.Sleep(3 * time.Second)

	csi := CommentServiceImpl{}
	id1, err := csi.AddCommentByAll(4, 1, "good video", time.Now())
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	_, err = csi.AddCommentByAll(4, 2, "good one", time.Now())
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	id3, err := csi.AddCommentByAll(14, 2, "hello you", time.Now())
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	_, err = csi.AddCommentByAll(4, 2, "haha", time.Now())
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	_, err = csi.AddCommentByAll(14, 1, "yes", time.Now())
	assert.Nil(t, err)

	cnt, err := csi.CountCommentsByToVideoId(1)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(2))
	cnt, err = csi.CountCommentsByToVideoId(2)
	assert.Nil(t, err)
	assert.Equal(t, cnt, int64(3))

	err = csi.DeleteCommentById(id1)
	assert.Nil(t, err)
	err = csi.DeleteCommentById(id3)
	assert.Nil(t, err)
	//然后查数据库对不对

	<-forever
}
