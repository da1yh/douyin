package rabbitmq

import (
	"context"
	"douyin/dao"
	"douyin/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"strings"
	"time"
)

type CommentMQ struct {
	RabbitMQ
}

var Commentmq *CommentMQ

func InitCommentMQ() {
	Commentmq = &CommentMQ{}
	Commentmq.queName = "commentmq"
	ch, err := Conn.Channel()
	failOnError(err, "failed to open a channel")
	Commentmq.channel = ch

	go Commentmq.Consume()

}

func (mq CommentMQ) Produce(tp, fromUserId, toVideoId, content, createDate, commentId string) {

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
			Body:        []byte(tp + util.MQSplit + fromUserId + util.MQSplit + toVideoId + util.MQSplit + content + util.MQSplit + createDate + util.MQSplit + commentId),
		})
	failOnError(err, "failed to publish a message")
}

func (mq CommentMQ) Consume() {
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
			tp, fromUserIdStr, toVideoIdStr, content, createDateStr, commentIdStr := arr[0], arr[1], arr[2], arr[3], arr[4], arr[5]
			if tp == util.MQCancelCommentType {
				commentId, _ := strconv.ParseInt(commentIdStr, 10, 64)
				err := dao.DeleteCommentById(commentId)
				failOnError(err, "failed to delete comment info to database")
			} else {
				fromUserId, _ := strconv.ParseInt(fromUserIdStr, 10, 64)
				toVideoId, _ := strconv.ParseInt(toVideoIdStr, 10, 64)
				createDate, _ := time.Parse(util.TimeFormat, createDateStr)
				_, err := dao.AddCommentByAll(fromUserId, toVideoId, content, createDate)
				failOnError(err, "failed to add comment info to database")
			}
		}
	}()

	<-forever

}
