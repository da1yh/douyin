package rabbitmq

import (
	"douyin/dao"
	"douyin/util"
	"testing"
	"time"
)

func TestRmq(t *testing.T) {
	dao.InitDb()
	InitRabbitMQ()
	forever := make(chan struct{})
	InitFavoriteMQ()
	time.Sleep(3 * time.Second)
	Favoritemq.Produce(util.MQLikeType, "4", "2")
	Favoritemq.Produce(util.MQLikeType, "14", "2")
	Favoritemq.Produce(util.MQLikeType, "4", "3")
	Favoritemq.Produce(util.MQDisLikeType, "4", "2")
	Favoritemq.Produce(util.MQDisLikeType, "4", "3")
	<-forever
}

func Test1(t *testing.T) {
	InitRabbitMQ()
	forever := make(chan struct{})
	InitTest1mq()
	InitTest2mq()
	//assert.NotNil(t, Test1mq.channel)
	time.Sleep(3 * time.Second)
	//go func() {
	// 用协程或者不用协程
	//go Test1mq.Produce("hello")
	//go Test2mq.Produce("niko")
	//go Test2mq.Produce("m0nesy")
	//go Test1mq.Produce("world")
	//go Test2mq.Produce("hooxi")
	//go Test1mq.Produce("major")
	//go Test1mq.Produce("champion")
	//go Test2mq.Produce("jks")
	//go Test1mq.Produce("last")

	Test1mq.Produce("hello")
	Test2mq.Produce("niko")
	Test2mq.Produce("m0nesy")
	Test1mq.Produce("world")
	Test2mq.Produce("hooxi")
	Test1mq.Produce("major")
	Test1mq.Produce("champion")
	Test2mq.Produce("jks")
	Test1mq.Produce("last")
	//}()
	<-forever
}

//func Test2(t *testing.T) {
//	//go func() {
//	assert.NotNil(t, Test1mq)
//	Test1mq.Produce("hello")
//	Test1mq.Produce("world")
//	Test1mq.Produce("major")
//	Test1mq.Produce("champion")
//	Test1mq.Produce("last")
//	//}()
//}
