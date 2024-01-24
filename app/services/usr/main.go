package main

import (
	"errors"
	"fmt"

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
	defer log.Sync()
	if err != nil {
		fmt.Println(err)
		return
	}

	run(log)
}

func run(log *zap.SugaredLogger) error {
	var cfg Config

	err := envconfig.Process("usr", &cfg)
	if err != nil {
		panic(err)
	}

	log.Infow("startup", "STATUS", "OK")
	log.Infow("usr config", "config debug value:", cfg.Debug)

	return nil
}

// TODO different configs and change logLevel
func initLogger(service string) (*zap.SugaredLogger, error) {
	var cfg LogConfig

	err := envconfig.Process("usr", &cfg)
	if err != nil {
		return nil, err
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
