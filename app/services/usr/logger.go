package main

import (
	"errors"

	"go.uber.org/zap"
)

type LogConfig struct {
	Mode string `default:"debug"`
}

func InitLogger(service string, cfg *LogConfig) (*zap.SugaredLogger, error) {
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
