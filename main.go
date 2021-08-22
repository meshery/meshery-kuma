package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"

	// "github.com/layer5io/meshkit/tracing"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshery-kuma/kuma"
	"github.com/layer5io/meshery-kuma/kuma/oam"
)

var (
	serviceName = "kuma-adaptor"
	version     = "none"
	gitsha      = "none"
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
	cfg, err := config.New(configprovider.ViperKey)
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

	kubeconfigHandler, err := config.NewKubeconfigBuilder(configprovider.ViperKey)
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
	service.Version = version
	service.GitSHA = gitsha

	go registerCapabilities(service.Port, log)

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

	err = os.Setenv("KUBECONFIG", fmt.Sprintf("%s/.meshery/kubeconfig.yaml", utils.GetHome()))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func mesheryServerAddress() string {
	meshReg := os.Getenv("MESHERY_SERVER")

	if meshReg != "" {
		if strings.HasPrefix(meshReg, "http") {
			return meshReg
		}

		return "http://" + meshReg
	}

	return "http://localhost:9081"
}

func serviceAddress() string {
	svcAddr := os.Getenv("SERVICE_ADDR")

	if svcAddr != "" {
		return svcAddr
	}

	return "mesherylocal.layer5.io"
}

func registerCapabilities(port string, log logger.Handler) {
	// Register workloads
	if err := oam.RegisterWorkloads(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Info(err.Error())
	}

	// Register traits
	if err := oam.RegisterTraits(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Info(err.Error())
	}
}
