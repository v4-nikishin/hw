package publisher

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"gopkg.in/square/go-jose.v2/json"
)

type Publisher struct {
	log             *logger.Logger
	publishedEvents map[string]struct{}

	uri          string // AMQP URI
	exchange     string // Durable AMQP exchange name
	exchangeType string // Exchange type - direct|fanout|topic|x-custom
	routingKey   string // AMQP routing key
	body         string // Body of message
	reliable     bool   // Wait for the publisher confirmation before exiting
}

func New(logger *logger.Logger) *Publisher {
	return &Publisher{
		log:             logger,
		publishedEvents: make(map[string]struct{}),

		uri:          "amqp://guest:guest@localhost:5672/",
		exchange:     "calendar-exchange",
		exchangeType: "direct",
		routingKey:   "calendar-key",
		reliable:     true,
	}
}

func (p *Publisher) Publish(e *pb.Events) {
	const format = "2006-01-02 15:04:00"
	for _, evt := range e.Events {
		beginStr := evt.GetDate() + " " + evt.GetBegin()
		begin, err := time.Parse(format, beginStr)
		if err != nil {
			p.log.Error("failed to parse date: " + beginStr)
			continue
		}
		endStr := evt.GetDate() + " " + evt.GetEnd()
		end, err := time.Parse(format, endStr)
		if err != nil {
			p.log.Error("failed to parse date: " + endStr)
			continue
		}
		if _, ok := p.publishedEvents[evt.Uuid]; ok {
			continue
		}
		now := time.Now().UTC()
		if (now.After(begin) || now.Equal(begin)) && (now.Before(end) || now.Equal(end)) {
			fmt.Printf("Event: %v\n", evt)
			buf, err := json.Marshal(evt)
			if err != nil {
				p.log.Error("failed to marshal event: " + err.Error())
				continue
			}
			p.body = string(buf)
			p.publish()
			if err != nil {
				p.log.Error("failed to publish event: " + err.Error())
				continue
			}
			p.publishedEvents[evt.Uuid] = struct{}{}
		}
	}
}

func (p *Publisher) publish() error {
	// This function dials, connects, declares, publishes, and tears down,
	// all in one go. In a real service, you probably want to maintain a
	// long-lived connection as state, and publish against that.

	p.log.Info(fmt.Sprintf("dialing %q", p.uri))
	connection, err := amqp.Dial(p.uri)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer connection.Close()

	p.log.Info("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	p.log.Info(fmt.Sprintf("got Channel, declaring %q Exchange (%q)", p.exchangeType, p.exchange))
	if err := channel.ExchangeDeclare(
		p.exchange,     // name
		p.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("exchange Declare: %w", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if p.reliable {
		p.log.Info("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %w", err)
		}

		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer p.confirmOne(confirms)
	}

	p.log.Info(fmt.Sprintf("declared Exchange, publishing %dB body (%q)", len(p.body), p.body))
	if err = channel.Publish(
		p.exchange,   // publish to an exchange
		p.routingKey, // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(p.body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %w", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func (p *Publisher) confirmOne(confirms <-chan amqp.Confirmation) {
	p.log.Info("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		p.log.Info(fmt.Sprintf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag))
	} else {
		p.log.Error(fmt.Sprintf("failed delivery of delivery tag: %d", confirmed.DeliveryTag))
	}
}
