package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/publisher"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	cfg, err := config.LoadSchedulerConfig(configFile)
	if err != nil {
		fmt.Printf("failed to configure service %s\n", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger, os.Stdout)

	addr := net.JoinHostPort(cfg.Server.GRPC.Host, cfg.Server.GRPC.Port)
	logg.Info("trying to connect to grpc server on " + addr)

	ctxT, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	conn, err := grpc.DialContext(ctxT, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		logg.Error("failed to dial to grpc server: " + err.Error())
		os.Exit(1)
	}
	client := pb.NewCalendarClient(conn)

	p, err := publisher.New(cfg.Publisher, logg)
	if err != nil {
		logg.Error("failed to create publisher: " + err.Error())
		os.Exit(1)
	}
	defer p.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	ticker := time.NewTicker(time.Duration(cfg.Scheduler.CheckPeriod) * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				e, err := client.GetEvents(ctx, &emptypb.Empty{})
				if err != nil {
					logg.Error("failed to get events")
					cancel()
					return
				}
				for _, evt := range e.Events {
					p.Publish(evt)
				}
			}
		}
	}()

	logg.Info("scheduler is running...")

	<-ctx.Done()

	logg.Info("...scheduler is stopped")
}
