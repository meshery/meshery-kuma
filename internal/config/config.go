package config

import (
	"path"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/status"
	configprovider "github.com/layer5io/meshkit/config/provider"
	"github.com/layer5io/meshkit/utils"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

const (
	Development = "development"
	Production  = "production"
	// OAM Metadata constants
	OAMAdapterNameMetadataKey       = "adapter.meshery.io/name"
	OAMComponentCategoryMetadataKey = "ui.meshery.io/category"
)

var (
	configRootPath = path.Join(utils.GetHome(), ".meshery")
	KumaOperation  = strings.ToLower(smp.ServiceMesh_KUMA.Enum().String())

	ServerConfig = map[string]string{
		"name":     smp.ServiceMesh_KUMA.Enum().String(),
		"type":     "adapter",
		"port":     "10007",
		"traceurl": status.None,
	}

	MeshSpec = map[string]string{
		"name":    smp.ServiceMesh_KUMA.Enum().String(),
		"status":  status.NotInstalled,
		"version": status.None,
	}

	Operations = getOperations(common.Operations)
)

// New creates a new config instance
func New(provider string) (h config.Handler, err error) {
	opts := configprovider.Options{
		FilePath: configRootPath,
		FileName: "kuma",
		FileType: "yaml",
	}

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		h, err = configprovider.NewViper(opts)
		if err != nil {
			return nil, err
		}
	case configprovider.InMemKey:
		h, err = configprovider.NewInMem(opts)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrEmptyConfig
	}

	// Setup server config
	if err := h.SetObject(adapter.ServerKey, ServerConfig); err != nil {
		return nil, err
	}

	// Setup mesh config
	if err := h.SetObject(adapter.MeshSpecKey, MeshSpec); err != nil {
		return nil, err
	}

	// Setup Operations Config
	if err := h.SetObject(adapter.OperationsKey, Operations); err != nil {
		return nil, err
	}

	return h, nil
}

func NewKubeconfigBuilder(provider string) (config.Handler, error) {
	opts := configprovider.Options{
		FilePath: configRootPath,
		FileType: "yaml",
		FileName: "kubeconfig",
	}

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}
	return nil, config.ErrEmptyConfig
}

// RootPath returns the configRootPath
func RootPath() string {
	return configRootPath
}
