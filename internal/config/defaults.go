package config

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
		InstallKumav071: map[string]interface{}{
			"type": "0",
			"properties": map[string]string{
				"description": "Kuma service mesh (0.7.1)",
				"version":     "0.7.1",
			},
		},
		InstallKumav070: map[string]interface{}{
			"type": "0",
			"properties": map[string]string{
				"description": "Kuma service mesh (0.7.0)",
				"version":     "0.7.0",
			},
		},
		InstallKumav060: map[string]interface{}{
			"type": "0",
			"properties": map[string]string{
				"description": "Kuma service mesh (0.6.0)",
				"version":     "0.6.0",
			},
		},
		InstallSampleBookInfo: map[string]interface{}{
			"type": "1",
			"properties": map[string]string{
				"description": "BookInfo",
				"version":     "latest",
			},
		},
		ValidateSmiConformance: map[string]interface{}{
			"type": "3",
			"properties": map[string]string{
				"description": "SMI conformance test",
				"version":     "latest",
			},
		},
	}

	// Viper configuration
	filepath = "/root/.kuma"
	filename = "config.yml"
	filetype = "yaml"
)
