package main

import (
	"log"

	"github.com/gromnsk/money.io/pkg/config"
	"github.com/gromnsk/money.io/pkg/service"
)

func main() {
	// Load ENV configuration
	cfg := new(config.Config)
	if err := cfg.Load(config.SERVICENAME); err != nil {
		log.Fatal(err)
	}

	service.Run(cfg)
}
