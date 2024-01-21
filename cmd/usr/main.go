package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug bool
}

func main() {
	var cfg Config
	err := envconfig.Process("usr", &cfg)
	if err != nil {
		panic(err)
	}
	println("debug:", cfg.Debug)
}
