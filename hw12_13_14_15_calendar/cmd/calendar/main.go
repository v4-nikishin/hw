package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	internalhttp "github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage/sql"
	"google.golang.org/grpc"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("failed to configure service %s\n", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger, os.Stdout)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var repo app.Storage
	if cfg.DB.Type == string(storage.DBTypeSQL) {
		if repo, err = sqlstorage.New(ctx, cfg.DB.SQL, logg); err != nil {
			logg.Error("failed to create sql storage: " + err.Error())
			return
		}
	} else {
		repo = memorystorage.New()
	}
	defer repo.Close()

	calendar := app.New(logg, repo)

	var server app.Server
	if cfg.Server.Proto == "grpc" {
		s := grpc.NewServer(grpc.ChainUnaryInterceptor())
		server = internalgrpc.NewServer(cfg.Server, logg, s, calendar)
		pb.RegisterCalendarServer(s, server.(pb.CalendarServer))
	} else {
		server = internalhttp.NewServer(cfg.Server, logg, calendar)
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop " + cfg.Server.Proto + " server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start " + cfg.Server.Proto + " server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
