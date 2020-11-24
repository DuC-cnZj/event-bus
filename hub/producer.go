package hub

import (
	"github.com/streadway/amqp"
)

type ProducerInterface interface {
	GetConn() *amqp.Connection
	GetChannel() *amqp.Channel
	GetQueueName() string
	GetKind() string
	Publish(Message) error
	Close()
}

type ProducerBase struct {
	pm        ProducerManagerInterface
	queueName string
	queue     amqp.Queue
	conn      *amqp.Connection
	channel   *amqp.Channel
	exchange  string

	kind string
	hub  Interface

	closed atomicBool

	closeChan chan *amqp.Error
}

var DefaultExchange = "duc_exchange"
