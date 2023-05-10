package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Test2MQ struct {
	RabbitMQ
}

var Test2mq *Test2MQ

//var message string = "test for one"

func InitTest2mq() {
	Test2mq = &Test2MQ{}
	Test2mq.queName = "test2mq"
	ch, err := Conn.Channel()
	failOnError(err, "failed to open a channel")
	Test2mq.channel = ch

	//Test1mq.queue = &q

	go Test2mq.Consume()

}

func (mq Test2MQ) Produce(msg string) {

	_, err := mq.channel.QueueDeclare(
		Test2mq.queName, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
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
	log.Printf(" [2] Sent %s\n", msg)
}

func (mq Test2MQ) Consume() {
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
			log.Printf("[2] received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] waiting for messages. To exit press CTRL+C")

	<-forever

}
