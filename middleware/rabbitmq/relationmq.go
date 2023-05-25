package rabbitmq

import (
	"context"
	"douyin/dao"
	"douyin/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"strings"
)

type RelationMQ struct {
	RabbitMQ
}

var Relationmq *RelationMQ

func InitRelationMQ() {
	Relationmq = &RelationMQ{}
	Relationmq.queName = "relationmq"
	ch, err := Conn.Channel()
	failOnError(err, "failed to open a channel")
	Relationmq.channel = ch

	go Relationmq.Consume()

}

func (mq RelationMQ) Produce(tp, fromUserId, toUserID string) {

	_, err := mq.channel.QueueDeclare(
		mq.queName, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "failed to declare a queue")

	ctx := context.Background()
	err = mq.channel.PublishWithContext(ctx,
		"",         // exchange
		mq.queName, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(tp + util.MQSplit + fromUserId + util.MQSplit + toUserID),
		})
	failOnError(err, "failed to publish a message")
}

func (mq RelationMQ) Consume() {
	_, err := mq.channel.QueueDeclare(
		mq.queName, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	failOnError(err, "failed to declare a queue")

	msgs, err := mq.channel.Consume(
		mq.queName, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			arr := strings.Split(string(d.Body), util.MQSplit)
			tp, fromUserIdStr, toUserIdStr := arr[0], arr[1], arr[2]
			fromUserId, _ := strconv.ParseInt(fromUserIdStr, 10, 64)
			toUserId, _ := strconv.ParseInt(toUserIdStr, 10, 64)
			res, _ := dao.CheckRelationByBothId(fromUserId, toUserId)
			if tp == util.MQFollowType {
				if !res {
					err = dao.AddRelationByBothId(fromUserId, toUserId)
					if err != nil {
						failOnError(err, "sync to database failed")
					}
				}
			} else {
				if res {
					err = dao.DeleteRelationByBothId(fromUserId, toUserId)
					if err != nil {
						failOnError(err, "sync to database failed")
					}
				}
			}
		}
	}()

	<-forever

}
