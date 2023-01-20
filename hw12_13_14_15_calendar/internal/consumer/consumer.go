package consumer

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
)

type Consumer struct {
	log *logger.Logger

	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error

	uri          string        // AMQP URI
	exchange     string        // Durable, non-auto-deleted AMQP exchange name
	exchangeType string        // Exchange type - direct|fanout|topic|x-custom
	queue        string        // Ephemeral AMQP queue name
	bindingKey   string        // AMQP binding key
	consumerTag  string        // AMQP consumer tag (should not be blank)
	lifetime     time.Duration // lifetime of process before shutdown (0s=infinite)
}

func NewConsumer(cfg config.ConsumerConf, log *logger.Logger) *Consumer {
	return &Consumer{
		log: log,

		conn:    nil,
		channel: nil,
		done:    make(chan error),

		uri:          cfg.URI,
		exchange:     cfg.Exchange,
		exchangeType: cfg.ExchangeType,
		queue:        cfg.Queue,
		bindingKey:   cfg.BindingKey,
		consumerTag:  cfg.ConsumerTag,
		lifetime:     0,
	}
}

func (c *Consumer) Consume() error {
	var err error

	c.log.Info(fmt.Sprintf("dialing %q", c.uri))
	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	go func() {
		c.log.Info(fmt.Sprintf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error))))
	}()

	c.log.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	c.log.Info(fmt.Sprintf("got Channel, declaring Exchange (%q)", c.exchange))
	if err = c.channel.ExchangeDeclare(
		c.exchange,     // name of the exchange
		c.exchangeType, // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("exchange Declare: %w", err)
	}

	c.log.Info(fmt.Sprintf("declared Exchange, declaring Queue %q", c.queue))
	queue, err := c.channel.QueueDeclare(
		c.queue, // name of the queue
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // noWait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %w", err)
	}

	c.log.Info(fmt.Sprintf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, c.bindingKey))

	if err = c.channel.QueueBind(
		queue.Name,   // name of the queue
		c.bindingKey, // bindingKey
		c.exchange,   // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("queue Bind: %w", err)
	}

	c.log.Info(fmt.Sprintf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.consumerTag))
	deliveries, err := c.channel.Consume(
		queue.Name,    // name
		c.consumerTag, // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		return fmt.Errorf("queue consume: %w", err)
	}

	go c.handle(deliveries, c.done)

	return nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.consumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %w", err)
	}

	defer c.log.Info("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func (c *Consumer) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		c.log.Info(fmt.Sprintf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		))
		d.Ack(false)
	}
	c.log.Info("handle: deliveries channel closed")
	done <- nil
}
