package messaging

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type IMessageClient interface {
	ConnectToBroker(connectionStr string) error
	PublishToQueue(data []byte, queueName string) error
	SubscribeToQueue(queueName string, handlerFunc func(amqp.Delivery)) error

	Close()
}

type MessageClient struct {
	conn *amqp.Connection
}

func (m *MessageClient) ConnectToBroker(connectionStr string) error {
	if connectionStr == "" {
		panic("the connection str mustnt be null")
	}
	var err error
	m.conn, err = amqp.Dial(connectionStr)
	return err
}

func (m *MessageClient) PublishToQueue(body []byte, queueName string) error {
	if m.conn == nil {
		panic("before publish you must connect the RabbitMQ first")
	}

	ch, err := m.conn.Channel()
	defer ch.Close()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	failOnError(err, "Failed to publish a message")

	return nil
}

func (m *MessageClient) SubscribeToQueue(queueName string, handlerFunc func(amqp.Delivery)) error {
	ch, err := m.conn.Channel()
	//defer ch.Close()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer")
	go consumeLoop(msgs, handlerFunc)
	return nil
}

func (m *MessageClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

func consumeLoop(msgs <-chan amqp.Delivery, handlerFunc func(amqp.Delivery)) {
	for d := range msgs {
		handlerFunc(d)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
		panic(fmt.Sprintf("%s: %s", err, msg))
	}
}
