package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func ParseConfig(process string, cfg interface{}) error {
	err := envconfig.Process(process, cfg)
	if err != nil {
		return err
	}
	fmt.Printf("Config parsed successful: %v\n", cfg)
	return nil
}
