package internalgrpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	memorystorage "github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAPI(t *testing.T) {
	event := pb.Event{
		Uuid:  "123e4567-e89b-12d3-a456-426655440002",
		Title: "Event title",
		User:  "Event user",
		Date:  "2023-01-18",
		Begin: "18:54:00",
		End:   "19:00:00",
	}

	logg := logger.New(config.LoggerConf{Level: "debug"}, os.Stdout)
	repo := memorystorage.New()
	calendar := app.New(logg, repo)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor())
	service := NewServer(config.ServerGRPC{Host: "localhost", Port: "50051"}, logg, server, calendar)
	pb.RegisterCalendarServer(server, service)

	addr := net.JoinHostPort(service.cfg.Host, service.cfg.Port)
	lsn, err := net.Listen("tcp", addr)
	require.NoError(t, err)

	log.Printf("starting server on %s", lsn.Addr().String())
	go func() {
		err := server.Serve(lsn)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := pb.NewCalendarClient(conn)

	t.Run("nil event", func(t *testing.T) {
		_, err := client.CreateEvent(context.Background(), nil)
		require.Error(t, err)
		respErr, _ := status.FromError(err)
		require.Equal(t, respErr.Code(), codes.Internal)
	})
	t.Run("create event", func(t *testing.T) {
		_, err := client.CreateEvent(context.Background(), &event)
		require.NoError(t, err)
	})
	t.Run("get event", func(t *testing.T) {
		e, err := client.GetEvent(context.Background(), &pb.EventId{Uuid: event.GetUuid()})
		require.NoError(t, err)
		require.Equal(t, e.GetUuid(), event.GetUuid())
	})
	t.Run("get invalid event", func(t *testing.T) {
		_, err := client.GetEvent(context.Background(), &pb.EventId{Uuid: "000"})
		require.Error(t, err)
		respErr, _ := status.FromError(err)
		require.Equal(t, respErr.Code(), codes.Internal)
	})
	t.Run("update event", func(t *testing.T) {
		event.Title = "Updated event title"
		_, err := client.UpdateEvent(context.Background(), &event)
		require.NoError(t, err)
	})
	t.Run("get updated event", func(t *testing.T) {
		e, err := client.GetEvent(context.Background(), &pb.EventId{Uuid: event.GetUuid()})
		require.NoError(t, err)
		require.Equal(t, e.GetUuid(), event.GetUuid())
		require.Equal(t, e.GetTitle(), "Updated event title")
	})
	t.Run("get event list", func(t *testing.T) {
		e, err := client.GetEvents(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)
		require.Equal(t, len(e.GetEvents()), 1)
	})
	t.Run("get events on date", func(t *testing.T) {
		d := pb.Date{Date: event.Date}
		e, err := client.GetEventsOnDate(context.Background(), &d)
		require.NoError(t, err)
		require.Equal(t, len(e.GetEvents()), 1)
	})
	t.Run("delete event", func(t *testing.T) {
		_, err := client.DeleteEvent(context.Background(), &pb.EventId{Uuid: event.GetUuid()})
		require.NoError(t, err)
	})
	t.Run("get deleted event", func(t *testing.T) {
		_, err := client.GetEvent(context.Background(), &pb.EventId{Uuid: event.GetUuid()})
		require.Error(t, err)
	})
	t.Run("get event list after deletion", func(t *testing.T) {
		e, err := client.GetEvents(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)
		require.Equal(t, len(e.GetEvents()), 0)
	})
}
