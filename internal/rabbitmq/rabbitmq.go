package rabbitmq

import (
	"fmt"

	"github.com/Dendyator/calendar/internal/logger" //nolint
	"github.com/streadway/amqp"                     //nolint
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logg    *logger.Logger
}

func New(dsn string, logg *logger.Logger) (*Client, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	logg.Info("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	logg.Info("RabbitMQ channel opened")

	return &Client{
		conn:    conn,
		channel: ch,
		logg:    logg,
	}, nil
}

func (c *Client) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		c.logg.Errorf("Failed to declare RabbitMQ queue: %s", err)
	} else {
		c.logg.Infof("RabbitMQ queue %s declared", name)
	}
	return err
}

func (c *Client) Publish(queue string, body []byte) error {
	err := c.channel.Publish(
		"",    // exchange
		queue, // routing key (queue name)
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		c.logg.Errorf("Failed to publish message to queue %s: %s", queue, err)
	} else {
		c.logg.Infof("Message published to queue %s", queue)
	}
	return err
}

func (c *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queue,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (c *Client) Close() {
	c.channel.Close()
	c.conn.Close()
	c.logg.Info("Closed RabbitMQ connection and channel")
}
