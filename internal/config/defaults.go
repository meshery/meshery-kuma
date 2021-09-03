package config

import (
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/status"
	configprovider "github.com/layer5io/meshkit/config/provider"
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
