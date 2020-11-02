package config

import (
	"github.com/layer5io/meshery-adapter-library/config"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
)

// New creates a new config instance
func New(provider string, environment string) (config.Handler, error) {

	// Default config
	opts = config.DefaultOpts

	// Config environment
	switch environment {
	case Production:
		opts = ProductionConfig
	case Development:
		opts = developmentConfig
	}

	// Config provider
	switch provider {
	case configprovider.Viper:
		return configprovider.NewViper(opts)
	case configprovider.InMem:
		return configprovider.NewInMem(opts)
	}

	return nil, ErrEmptyConfig
}
