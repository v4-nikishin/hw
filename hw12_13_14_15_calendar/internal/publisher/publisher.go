package publisher

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"gopkg.in/square/go-jose.v2/json"
)

type Publisher struct {
	log             *logger.Logger
	publishedEvents map[string]struct{}
	connection      *amqp.Connection
	channel         *amqp.Channel

	uri          string // AMQP URI
	exchange     string // Durable AMQP exchange name
	exchangeType string // Exchange type - direct|fanout|topic|x-custom
	routingKey   string // AMQP routing key
	body         string // Body of message
	reliable     bool   // Wait for the publisher confirmation before exiting
}

func New(cfg config.PublisherConf, logger *logger.Logger) (*Publisher, error) {
	p := &Publisher{
		log:             logger,
		publishedEvents: make(map[string]struct{}),

		uri:          cfg.URI,
		exchange:     cfg.Exchange,
		exchangeType: cfg.ExchangeType,
		routingKey:   cfg.RoutingKey,
		reliable:     cfg.Reliable,
	}

	var err error

	p.log.Info(fmt.Sprintf("dialing %q", p.uri))

	tryNum := 60
	for i := 0; i < tryNum; i++ {
		p.connection, err = amqp.Dial(p.uri)
		if err != nil {
			p.log.Info("Dialing...")
			time.Sleep(time.Second)
			continue
		}
	}
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	p.log.Info("got Connection, getting Channel")
	p.channel, err = p.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %w", err)
	}
	return p, nil
}

func (p *Publisher) Publish(evt *pb.Event) {
	const format = "2006-01-02 15:04:05"
	beginStr := fmt.Sprintf("%s %s", evt.GetDate(), evt.GetBegin())
	begin, err := time.Parse(format, beginStr)
	if err != nil {
		p.log.Error("failed to parse date: " + beginStr)
		return
	}
	endStr := fmt.Sprintf("%s %s", evt.GetDate(), evt.GetEnd())
	end, err := time.Parse(format, endStr)
	if err != nil {
		p.log.Error("failed to parse date: " + endStr)
		return
	}
	if _, ok := p.publishedEvents[evt.Uuid]; ok {
		return
	}
	now := time.Now().UTC()
	if (now.After(begin) || now.Equal(begin)) && (now.Before(end) || now.Equal(end)) {
		fmt.Printf("Event: %v\n", evt)
		buf, err := json.Marshal(evt)
		if err != nil {
			p.log.Error("failed to marshal event: " + err.Error())
			return
		}
		p.body = string(buf)
		err = p.publish()
		if err != nil {
			p.log.Error("failed to publish event: " + err.Error())
			return
		}
		p.publishedEvents[evt.Uuid] = struct{}{}
	}
}

func (p *Publisher) publish() error {
	p.log.Info(fmt.Sprintf("got Channel, declaring %q Exchange (%q)", p.exchangeType, p.exchange))
	if err := p.channel.ExchangeDeclare(
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
		if err := p.channel.Confirm(false); err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %w", err)
		}

		confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer p.confirmOne(confirms)
	}

	p.log.Info(fmt.Sprintf("declared Exchange, publishing %dB body (%q)", len(p.body), p.body))
	if err := p.channel.Publish(
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

func (p *Publisher) Close() {
	if err := p.connection.Close(); err != nil {
		p.log.Error("failed to close connection: " + err.Error())
	}
}
