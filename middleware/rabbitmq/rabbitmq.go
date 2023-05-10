package rabbitmq

import (
	"douyin/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQ struct {
	channel *amqp.Channel
	queue   *amqp.Queue
	queName string
}

var Conn *amqp.Connection

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func InitRabbitMQ() {
	conn, err := amqp.Dial(config.RabbitMQServerAddr)
	failOnError(err, "Failed to connect to RabbitMQ")
	Conn = conn
}
