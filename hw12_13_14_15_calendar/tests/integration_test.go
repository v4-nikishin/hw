package internalgrpc

import (
	"context"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/consumer"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSendEvent(t *testing.T) {
	f := "2006-01-02 15:04:05"

	now := time.Now().UTC()
	dateTime := now.Format(f)
	dateTime1 := now.Add(5 * time.Minute).Format(f)

	s := strings.Fields(dateTime)
	s1 := strings.Fields(dateTime1)

	uuid := (uuid.New()).String()
	event := pb.Event{
		Uuid:  uuid,
		Title: "Event title",
		User:  "Event user",
		Date:  s[0],
		Begin: s[1],
		End:   s1[1],
	}

	addr := net.JoinHostPort("", "50051")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := pb.NewCalendarClient(conn)

	logg := logger.New(config.LoggerConf{Level: "debug"}, os.Stdout)

	conf := config.ConsumerConf{
		URI:          "amqp://guest:guest@:5672/",
		Exchange:     "calendar-exchange",
		ExchangeType: "direct",
		Queue:        "calendar-queue",
		BindingKey:   "calendar-key",
		ConsumerTag:  "calendar-consumer",
		Lifetime:     10,
	}

	c, err := consumer.NewConsumer(conf, logg)
	require.NoError(t, err)
	defer func() {
		err = c.Shutdown()
		require.NoError(t, err)
	}()
	err = c.Consume()
	require.NoError(t, err)

	t.Run("check consume event", func(t *testing.T) {
		_, err := client.CreateEvent(context.Background(), &event)
		require.NoError(t, err)

		time.Sleep(time.Duration(conf.Lifetime * uint64(time.Second)))

		require.True(t, c.IsSentEvent(uuid))
	})
}
