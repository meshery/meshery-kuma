package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"

	// "github.com/layer5io/meshkit/tracing"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshery-kuma/kuma"
)

var (
	serviceName = "kuma-adaptor"
	environment = "development"
	provider    = configprovider.ViperKey
)

// main is the entrypoint of the adaptor
func main() {

	// Initialize Logger instance
	log, err := logger.New(serviceName, logger.Options{
		Format: logger.SyslogLogFormat,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize application specific configs and dependencies
	// App and request config
	cfg, err := config.New(provider, environment)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	service := &grpc.Service{}
	err = cfg.GetObject(adapter.ServerKey, service)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	kubeconfigHandler, err := config.NewKubeconfigBuilder(configprovider.ViperKey, environment)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Initialize Handler intance
	handler := kuma.New(cfg, log, kubeconfigHandler)
	handler = adapter.AddLogger(log, handler)

	service.Handler = handler
	service.Channel = make(chan interface{}, 10)
	service.StartedAt = time.Now()

	// Server Initialization
	log.Info("Adaptor Listening at port: ", service.Port)
	err = grpc.Start(service, nil)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	err := os.MkdirAll(path.Join(utils.GetHome(), ".meshery", "bin"), 0750)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
