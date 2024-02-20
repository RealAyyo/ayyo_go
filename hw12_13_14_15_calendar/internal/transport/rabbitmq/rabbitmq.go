package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	eventQueue = "calendar_event"
)

type RabbitMQ struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	EventsQueue amqp.Queue
}

func New(user string, password string, host string, port string) (*RabbitMQ, error) {
	uri := "amqp://" + user + ":" + password + "@" + host + ":" + port + "/"
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMQ) QueueDeclare() error {
	eventsQueue, err := r.ch.QueueDeclare(
		eventQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.EventsQueue = eventsQueue
	return nil
}

func (r *RabbitMQ) Send(ctx context.Context, queueName string, body []byte) error {
	err := r.ch.PublishWithContext(ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	return err
}

func (r *RabbitMQ) Consume(queueName string, handler func(message amqp.Delivery)) error {
	messages, err := r.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	go func() {
		for d := range messages {
			handler(d)
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() error {
	if r.conn.IsClosed() {
		return nil
	}
	return r.conn.Close()
}
