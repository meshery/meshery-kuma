package config

import (
	"github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/status"
)

const (
	Development = "development"
)

var (
	DevelopmentConfig = config.Options{
		ServerConfig:   developmentServerConfig,
		MeshSpec:       developmentMeshSpec,
		MeshInstance:   developmentMeshInstance,
		ProviderConfig: developmentProviderConfig,
		Operations:     developmentOperations,
	}

	developmentServerConfig = map[string]string{
		"name":    "kuma-adapter",
		"port":    "10007",
		"version": "v1.0.0",
	}

	developmentMeshSpec = map[string]string{
		"name":     "kuma",
		"status":   status.NotInstalled,
		"traceurl": "none",
		"version":  "none",
	}

	developmentMeshInstance = map[string]string{}

	developmentProviderConfig = map[string]string{}

	developmentOperations = map[string]interface{}{
		InstallKumav071: map[string]interface{}{
			"type": "0",
			"properties": map[string]string{
				"description": "Install Kuma service mesh (0.7.1)",
				"version":     "0.7.1",
			},
		},
		InstallKumav070: map[string]interface{}{
			"type": "0",
			"properties": map[string]string{
				"description": "Install Kuma service mesh (0.7.0)",
				"version":     "0.7.0",
			},
		},
		InstallKumav060: map[string]interface{}{
			"type": "0",
			"properties": map[string]string{
				"description": "Install Kuma service mesh (0.6.0)",
				"version":     "0.6.0",
			},
		},
		InstallSampleBookInfo: map[string]interface{}{
			"type": "1",
			"properties": map[string]string{
				"description": "Install BookInfo Application",
				"version":     "latest",
			},
		},
		ValidateSmiConformance: map[string]interface{}{
			"type": "3",
			"properties": map[string]string{
				"description": "Validate SMI conformance",
				"version":     "latest",
			},
		},
	}

	// Viper configuration
	filepath = "/root/.kuma"
	filename = "config.yml"
	filetype = "yaml"
)
