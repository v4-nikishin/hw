package internalgrpc

import (
	"context"
	"net"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedCalendarServer

	cfg    config.ServerConf
	log    *logger.Logger
	server *grpc.Server
	app    *app.App
}

func NewServer(cfg config.ServerConf, logger *logger.Logger, server *grpc.Server, app *app.App) *Server {
	return &Server{cfg: cfg, log: logger, server: server, app: app}
}

func (s *Server) convertToGrpcEvent(event *storage.Event) *pb.Event {
	return &pb.Event{
		Uuid:  event.UUID,
		Title: event.Title,
		User:  event.User,
		Date:  event.Date,
		Begin: event.Begin,
		End:   event.End,
	}
}

func (s *Server) convertToStorageEvent(e *pb.Event) *storage.Event {
	return &storage.Event{
		UUID:  e.GetUuid(),
		Title: e.GetTitle(),
		User:  e.GetUser(),
		Date:  e.GetDate(),
		Begin: e.GetBegin(),
		End:   e.GetEnd(),
	}
}

func (s *Server) CreateEvent(ctx context.Context, e *pb.Event) (*emptypb.Empty, error) {
	if e == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}
	evt := s.convertToStorageEvent(e)
	if err := s.app.CreateEvent(*evt); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetEvent(ctx context.Context, id *pb.EventId) (*pb.Event, error) {
	if id == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}
	event, err := s.app.GetEvent(id.GetUuid())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return s.convertToGrpcEvent(&event), nil
}

func (s *Server) UpdateEvent(ctx context.Context, e *pb.Event) (*emptypb.Empty, error) {
	if e == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}
	evt := s.convertToStorageEvent(e)
	if err := s.app.UpdateEvent(evt.UUID, *evt); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetEvents(context.Context, *emptypb.Empty) (*pb.Events, error) {
	storageEvts, err := s.app.Events()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	qrpcEventes := pb.Events{}

	//nolint:typecheck
	for _, e := range storageEvts {
		e := e
		qrpcEventes.Events = append(qrpcEventes.Events, s.convertToGrpcEvent(&e))
	}
	return &qrpcEventes, nil
}

func (s *Server) DeleteEvent(ctx context.Context, id *pb.EventId) (*emptypb.Empty, error) {
	if id == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}
	err := s.app.DeleteEvent(id.GetUuid())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.cfg.Host, s.cfg.Port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		s.log.Error("failed to start grpc server: " + err.Error())
		return err
	}
	s.log.Info("starting grpc server on " + lsn.Addr().String())
	if err := s.server.Serve(lsn); err != nil {
		s.log.Error("failed to start grpc server: " + err.Error())
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.GracefulStop()
	s.log.Info("...calendar is stopped")
	return nil
}
