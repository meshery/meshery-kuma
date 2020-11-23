package config

import (
	"path"

	"github.com/layer5io/meshery-adapter-library/config"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshkit/utils"
)

const (
	KumaOperation = "kuma"
	Development   = "development"
	Production    = "production"
)

var (
	configRootPath = path.Join(utils.GetHome(), ".meshery")
)

// New creates a new config instance
func New(provider string, environment string) (config.Handler, error) {

	opts := DevelopmentConfig

	// Config environment
	switch environment {
	case Production:
		opts = ProductionConfig
	case Development:
		opts = DevelopmentConfig
	}

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}

	return nil, ErrEmptyConfig
}

func NewKubeconfigBuilder(provider string, environment string) (config.Handler, error) {

	opts := configprovider.Options{}

	// Config environment
	switch environment {
	case Production:
		opts.ProviderConfig = productionKubeConfig
	case Development:
		opts.ProviderConfig = developmentKubeConfig
	}

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}
	return nil, ErrEmptyConfig
}
