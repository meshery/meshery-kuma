package config

import (
	"os/user"
)

var (

	// server holds server configuration
	server = map[string]string{
		"name":    "kuma-adapter",
		"port":    "10007",
		"version": "v1.0.0",
	}

	// mesh holds mesh configuration
	mesh = map[string]string{
		"name":     "kuma",
		"status":   "not installed",
		"traceurl": "http://localhost:14268/api/traces",
		"version":  "0.6.0",
	}

	// operations holds the supported operations inside mesh
	operations = map[string]interface{}{
		InstallKuma: map[string]interface{}{
			"type": "0",
			"properties": map[string]string{
				"description": "Install kuma service mesh",
				"version":     "latest",
			},
		},
		InstallSample: map[string]interface{}{
			"type": "1",
			"properties": map[string]string{
				"description": "Install sample application",
				"version":     "latest",
			},
		},
		RunSmiConformance: map[string]interface{}{
			"type": "3",
			"properties": map[string]string{
				"description": "Run SMI conformance test",
				"version":     "latest",
			},
		},
	}

	// Viper configuration
	filepath = "/root/.kuma"
	filename = "config.yml"
	filetype = "yaml"
)

// GetHome returns the home path
func GetHome() string {
	usr, _ := user.Current()
	return usr.HomeDir
}
