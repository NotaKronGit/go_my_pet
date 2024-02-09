package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	Debug bool
}

type LogConfig struct {
	Mode string
}

func main() {
	log, err := initLogger("service")
	defer log.Sync() // nolint: errcheck
	if err != nil {
		fmt.Println(err)
		return
	}

	run(log)
}

func run(log *zap.SugaredLogger) {
	var cfg Config

	err := envconfig.Process("usr", &cfg)
	if err != nil {
		panic(err)
	}

	log.Infow("startup", "STATUS", "OK")
	log.Infow("usr config", "config debug value:", cfg.Debug)

	ctx, stop := context.WithCancel(context.Background())
	go waitSignalExit(stop)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Infow("startup", "STATUS", "Break process")
				return
			case <-time.After(1 * time.Second):
				process()
			}
		}
	}()

	wg.Wait()
	log.Infow("startup", "STATUS", "Graceful shutdown.")
}

// TODO different configs and change logLevel
func initLogger(service string) (*zap.SugaredLogger, error) {
	var cfg LogConfig

	err := envconfig.Process("usr", &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Mode == "" {
		cfg.Mode = "debug"
	}

	switch cfg.Mode {
	case "debug":
		return zap.Must(zap.NewDevelopment()).Sugar().Named(service), nil
	case "prod":
		return zap.Must(zap.NewProduction()).Sugar().Named(service), nil
	case "newNop":
		return zap.NewNop().Sugar().Named(service), nil
	}
	return nil, errors.New("unknown logger mode")
}

func waitSignalExit(cancel func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	cancel()
}

func process() {
	for i := 0; i < 9; i++ {
		fmt.Printf("step: %d\n", i)
		time.Sleep(time.Second)
	}
}
