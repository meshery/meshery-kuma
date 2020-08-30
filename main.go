package main

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/layer5io/gokit/logger"
	"github.com/layer5io/meshery-kuma/api/grpc"
	"github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshery-kuma/internal/tracing"
	"github.com/layer5io/meshery-kuma/kuma"
)

var (
	serviceName    = "kuma-adaptor"
	configProvider = "local"
	smiHelmPath    = "https://github.com/kumarabd/learn-layer5/raw/kumarabd/feature/helm/charts/smi-conformance-0.1.0.tgz"
	smiNamespace   = "meshery"
	kubeConfigPath string
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

	user, err := user.Current()
	if err != nil {
		log.Err("Cannot get current user", err.Error())
		os.Exit(1)
	}
	kubeConfigPath = fmt.Sprintf("%s/.kube/config", user.HomeDir)

	cfg.SetKey("smi-helm-path", smiHelmPath)
	cfg.SetKey("kube-config-path", kubeConfigPath)
	cfg.SetKey("smi-namespace", smiNamespace)

	// Initialize Tracing instance
	tracer, err := tracing.New(service.Name, service.TraceURL)
	if err != nil {
		log.Err("Tracing Init Failed", err.Error())
		os.Exit(1)
	}

	// Initialize Handler intance
	handler := kuma.New(cfg, log)
	handler = kuma.AddLogger(log, handler)
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
