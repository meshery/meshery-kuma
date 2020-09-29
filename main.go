package main

import (
	"fmt"
	"os"
	"time"

	"github.com/layer5io/gokit/logger"
	// "github.com/layer5io/gokit/tracing"
	"github.com/layer5io/gokit/utils"
	"github.com/layer5io/meshery-kuma/api/grpc"
	"github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshery-kuma/kuma"
)

var (
	serviceName    = "kuma-adaptor"
	configProvider = "viper"
	kubeConfigPath = fmt.Sprintf("%s/.kube/config", utils.GetHome())
)

// main is the entrypoint of the adaptor
func main() {

	// Initialize Logger instance
	log, err := logger.New(serviceName)
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	// Initialize application specific configs and dependencies
	// App and request config
	cfg, err := config.New(configProvider)
	if err != nil {
		log.Err("Config Init Failed", err.Error())
		os.Exit(1)
	}
	service := &grpc.Service{}
	_ = cfg.Server(&service)
	cfg.SetKey("kube-config-path", kubeConfigPath)

	// // Initialize Tracing instance
	// tracer, err := tracing.New(service.Name, service.TraceURL)
	// if err != nil {
	// 	log.Err("Tracing Init Failed", err.Error())
	// 	os.Exit(1)
	// }

	// Initialize Handler intance
	handler := kuma.New(cfg, log)
	handler = kuma.AddLogger(log, handler)
	service.Handler = handler
	service.Channel = make(chan interface{}, 100)
	service.StartedAt = time.Now()

	// Server Initialization
	log.Info(fmt.Sprintf("Adaptor Started at: %s", service.Port))
	err = grpc.Start(service)
	if err != nil {
		log.Err("adapter crashed!!", err.Error())
		os.Exit(1)
	}
}
