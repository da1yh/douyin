package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Test1MQ struct {
	RabbitMQ
}

var Test1mq *Test1MQ

//var message string = "test for one"

func InitTest1mq() {
	Test1mq = &Test1MQ{}
	Test1mq.queName = "test1mq"
	ch, err := Conn.Channel()
	failOnError(err, "failed to open a channel")
	Test1mq.channel = ch

	//Test1mq.queue = &q

	go Test1mq.Consume()

}

func (mq Test1MQ) Produce(msg string) {

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
			Body:        []byte(msg),
		})
	failOnError(err, "failed to publish a message")
	log.Printf(" [1] Sent %s\n", msg)
}

func (mq Test1MQ) Consume() {
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
			log.Printf("[1] received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] waiting for messages. To exit press CTRL+C")

	<-forever

}
