package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/consumer"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/version"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/sender_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	cfg, err := config.LoadSenderConfig(configFile)
	if err != nil {
		fmt.Printf("failed to configure service %s\n", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger, os.Stdout)

	c, err := consumer.NewConsumer(cfg.Consumer, logg)
	if err != nil {
		logg.Error("failed to create consumer: " + err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := c.Shutdown(); err != nil {
			logg.Error("error during consumer shutdown: " + err.Error())
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		err := c.Consume()
		if err != nil {
			logg.Error("failed to consume: " + err.Error())
			cancel()
		}
	}()

	logg.Info("sender is running...")

	<-ctx.Done()

	logg.Info("...sender is stopped")
}
