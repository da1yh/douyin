package rabbitmq

import (
	"douyin/dao"
	"douyin/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddComment(t *testing.T) {
	dao.InitDb()
	_, err := dao.AddCommentByAll(4, 1, "nice video", time.Now())
	assert.Nil(t, err)
	time.Sleep(3 * time.Second)
	_, err = dao.AddCommentByAll(4, 2, "good one", time.Now())
	assert.Nil(t, err)
	time.Sleep(3 * time.Second)
	_, err = dao.AddCommentByAll(14, 2, "hello world", time.Now())
	assert.Nil(t, err)
}

// 实际上增加评论的数据库操作并不会用到mq，只有删除评论的操作会用到，所以这里produce只用删除评论
func TestCommentMQ(t *testing.T) {
	dao.InitDb()
	InitRabbitMQ()
	forever := make(chan struct{})
	InitCommentMQ()
	time.Sleep(3 * time.Second)
	Commentmq.Produce(util.MQCancelCommentType, util.MQEmptyValue, util.MQEmptyValue, util.MQEmptyValue, util.MQEmptyValue, "4")
	Commentmq.Produce(util.MQCancelCommentType, util.MQEmptyValue, util.MQEmptyValue, util.MQEmptyValue, util.MQEmptyValue, "5")
	<-forever
}
