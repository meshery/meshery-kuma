package main

import (
	"fmt"
	"os"
	"time"

	"github.com/layer5io/gokit/logger"
	"github.com/layer5io/meshery-kuma/api/grpc"
	"github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshery-kuma/internal/tracing"
	"github.com/layer5io/meshery-kuma/kuma"
)

var (
	configProvider = "local"
)

// main is the entrypoint of the adaptor
func main() {

	// Initialize application specific configs and dependencies
	// App and request config
	cfg, err := config.New(configProvider)
	if err != nil {
		fmt.Println("Config Init Failed", err.Error())
		os.Exit(1)
	}
	service := &grpc.Service{}
	_ = cfg.Server(&service)

	// Initialize Logger instance
	log, err := logger.New(service.Name)
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	// Initialize Tracing instance
	tracer, err := tracing.New(service.Name, service.TraceURL)
	if err != nil {
		fmt.Println("Tracing Init Failed", err.Error())
		os.Exit(1)
	}

	// Initialize Handler intance
	handler := kuma.New(cfg, log)
	service.Handler = handler
	service.Channel = make(chan interface{}, 100)
	service.StartedAt = time.Now()

	// Server Initialization
	log.Info(fmt.Sprintf("Adaptor Started at: %s", service.Port))
	err = grpc.Start(service, tracer)
	if err != nil {
		log.Err("Adaptor crashed!!", err.Error())
		os.Exit(1)
	}
}
