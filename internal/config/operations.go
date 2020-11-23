package config

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
)

var (
	ServiceName = "service_name"
)

func getOperations(dev adapter.Operations) adapter.Operations {

	dev[KumaOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "Kuma Service Mesh",
		Versions: []adapter.Version{
			"1.0.0",
			"0.7.3",
			"0.7.0",
			"0.6.0",
		},
		Templates: adapter.NoneTemplate,
		AdditionalProperties: map[string]string{
			ServiceName: KumaOperation,
		},
	}

	return dev
}
