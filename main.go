package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"

	// "github.com/layer5io/meshkit/tracing"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	"github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshery-kuma/kuma"
	"github.com/layer5io/meshery-kuma/kuma/oam"
	configprovider "github.com/layer5io/meshkit/config/provider"
	"github.com/layer5io/meshkit/utils/kubernetes"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	serviceName = "kuma-adapter"
	version     = "edge"
	gitsha      = "none"
)

// main is the entrypoint of the adapter
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

	go registerCapabilities(service.Port, log)        //Registering static capabilities
	go registerDynamicCapabilities(service.Port, log) //Registering latest capabilities periodically

	// Server Initialization
	log.Info("Adapter listening on port: ", service.Port)
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

	return "localhost"
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

func registerDynamicCapabilities(port string, log logger.Handler) {
	registerWorkloads(port, log)
	//Start the ticker
	const reRegisterAfter = 24
	ticker := time.NewTicker(reRegisterAfter * time.Hour)
	for {
		<-ticker.C
		registerWorkloads(port, log)
	}

}

func registerWorkloads(port string, log logger.Handler) {
	appVersion, chartVersion, err := getLatestValidAppVersionAndChartVersion()
	if err != nil {
		log.Info("Could not get latest service mesh version")
		return
	}
	log.Info("Registering latest workload components for version ", appVersion)
	// Register workloads
	if err := adapter.RegisterWorkLoadsDynamically(mesheryServerAddress(), serviceAddress()+":"+port, &adapter.DynamicComponentsConfig{
		TimeoutInMinutes: 10,
		URL:              "https://github.com/kumahq/charts/releases/download/kuma-" + chartVersion + "/kuma-" + chartVersion + ".tgz",
		GenerationMethod: adapter.HelmCHARTS,
		Config: manifests.Config{
			Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_KUMA)],
			MeshVersion: appVersion,
			Filter: manifests.CrdFilter{
				RootFilter:    []string{"$[?(@.kind==\"CustomResourceDefinition\")]"},
				NameFilter:    []string{"$..[\"spec\"][\"names\"][\"kind\"]"},
				VersionFilter: []string{"$[0]..spec.versions[0]"},
				GroupFilter:   []string{"$[0]..spec"},
				SpecFilter:    []string{"$[0]..openAPIV3Schema.properties.spec"},
				ItrFilter:     []string{"$[?(@.spec.names.kind"},
				ItrSpecFilter: []string{"$[?(@.spec.names.kind"},
				VField:        "name",
				GField:        "group",
			},
		},
		Operation: config.KumaOperation,
	}); err != nil {
		log.Info(err.Error())
		return
	}
	log.Info("Successfully registered latest service mesh components with Meshery Server at ", mesheryServerAddress())
}
func getLatestValidAppVersionAndChartVersion() (string, string, error) {
	release, err := utils.GetLatestReleaseTagsSorted("kumahq", "kuma")
	if err != nil {
		return "", "", kuma.ErrGetLatestRelease(err)
	}
	//loops through latest 10 app versions untill it finds one which is available in helm chart's index.yaml
	for i := range release {
		if chartVersion, err := kubernetes.HelmAppVersionToChartVersion("https://kumahq.github.io/charts", "kuma", release[len(release)-i-1]); err == nil {
			return release[len(release)-i-1], chartVersion, nil
		}

	}
	return "", "", kuma.ErrGetLatestRelease(err)
}
