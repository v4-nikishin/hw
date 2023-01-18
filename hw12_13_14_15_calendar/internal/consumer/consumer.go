package consumer

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
)

type Consumer struct {
	log *logger.Logger

	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error

	uri          string        //= flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     string        //= flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType string        //= flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        string        //= flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey   string        //= flag.String("key", "test-key", "AMQP binding key")
	consumerTag  string        //= flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	lifetime     time.Duration //= flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
}

func NewConsumer(log *logger.Logger) *Consumer {
	return &Consumer{
		log: log,

		conn:    nil,
		channel: nil,
		done:    make(chan error),

		uri:          "amqp://guest:guest@localhost:5672/",
		exchange:     "calendar-exchange",
		exchangeType: "direct",
		queue:        "calendar-queue",
		bindingKey:   "calendar-key",
		consumerTag:  "calendar-consumer",
		lifetime:     0,
	}
}
func (c *Consumer) Consume() error {
	var err error

	c.log.Info(fmt.Sprintf("dialing %q", c.uri))
	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	go func() {
		c.log.Info(fmt.Sprintf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error))))
	}()

	c.log.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
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
		return fmt.Errorf("exchange Declare: %s", err)
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
		return fmt.Errorf("queue Declare: %s", err)
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
		return fmt.Errorf("queue Bind: %s", err)
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
		return fmt.Errorf("queue consume: %s", err)
	}

	go c.handle(deliveries, c.done)

	return nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.consumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
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
