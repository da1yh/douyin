package rabbitmq

import (
	"douyin/dao"
	"douyin/util"
	"testing"
	"time"
)

func TestRelation(t *testing.T) {
	dao.InitDb()
	InitRabbitMQ()
	forever := make(chan struct{})
	InitRelationMQ()
	time.Sleep(3 * time.Second)
	// 测试关注
	// 4 -> 14
	// 4 -> 16
	// 16 -> 4 X
	// 14 -> 17
	//Relationmq.Produce(util.MQFollowType, "4", "14")
	//Relationmq.Produce(util.MQFollowType, "4", "16")
	//Relationmq.Produce(util.MQFollowType, "16", "4")
	//Relationmq.Produce(util.MQFollowType, "14", "17")

	Relationmq.Produce(util.MQUnfollowType, "16", "4")

	// 测试取关
	<-forever
}
