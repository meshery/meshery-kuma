package config

import (
	"github.com/layer5io/meshery-adapter-library/common"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-adapter-library/status"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	ServerDefaults = map[string]string{
		"name":     smp.ServiceMesh_KUMA.Enum().String(),
		"type":     "adapter",
		"port":     "10007",
		"traceurl": "none",
	}

	MeshSpecDefaults = map[string]string{
		"name":    smp.ServiceMesh_KUMA.Enum().String(),
		"status":  status.NotInstalled,
		"version": "none",
	}

	ProviderConfigDefaults = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "kuma",
	}

	KubeConfigDefaults = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "kubeconfig",
	}

	OperationsDefaults = getOperations(common.Operations)
)
