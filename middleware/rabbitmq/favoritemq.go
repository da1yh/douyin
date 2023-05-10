package rabbitmq

import (
	"context"
	"douyin/dao"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"strings"
)

type FavoriteMQ struct {
	RabbitMQ
}

var Favoritemq *FavoriteMQ

func InitFavoriteMQ() {
	Favoritemq = &FavoriteMQ{}
	Favoritemq.queName = "favoritemq"
	ch, err := Conn.Channel()
	failOnError(err, "failed to open a channel")
	Favoritemq.channel = ch

	go Favoritemq.Consume()

}

func (mq FavoriteMQ) Produce(tp, fromUserId, toVideoId string) {

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
			Body:        []byte(tp + "-" + fromUserId + "-" + toVideoId),
		})
	failOnError(err, "failed to publish a message")
}

func (mq FavoriteMQ) Consume() {
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
			arr := strings.Split(string(d.Body), "-")
			tp, fromUserIdStr, toVideoIdStr := arr[0], arr[1], arr[2]
			fromUserId, _ := strconv.ParseInt(fromUserIdStr, 10, 64)
			toVideoId, _ := strconv.ParseInt(toVideoIdStr, 10, 64)
			res, _ := dao.CheckFavoriteByBothId(fromUserId, toVideoId)
			if tp == "like" {
				if !res {
					err := dao.AddFavoriteByBothId(fromUserId, toVideoId)
					if err != nil {
						failOnError(err, "sync to database failed")
					}
				}
			} else {
				if res {
					err := dao.DeleteFavoriteByBothId(fromUserId, toVideoId)
					if err != nil {
						failOnError(err, "sync to database failed")
					}
				}
			}
		}
	}()

	<-forever

}
