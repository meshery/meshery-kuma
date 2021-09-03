package oam

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-kuma/internal/config"
)

var (
	basePath, _  = os.Getwd()
	workloadPath = filepath.Join(basePath, "templates", "oam", "workloads")
	traitPath    = filepath.Join(basePath, "templates", "oam", "traits")
)

type schemaDefinitionPathSet struct {
	oamDefinitionPath string
	jsonSchemaPath    string
	name              string
}

// RegisterWorkloads will register all of the workload definitions
// present in the path oam/workloads
//
// Registration process will send POST request to $runtime/api/oam/workload
func RegisterWorkloads(runtime, host string) error {
	oamRDP := []adapter.OAMRegistrantDefinitionPath{}

	pathSets, err := load(workloadPath)
	if err != nil {
		return err
	}

	for _, pathSet := range pathSets {
		metadata := map[string]string{
			config.OAMAdapterNameMetadataKey: config.KumaOperation,
		}

		if strings.HasSuffix(pathSet.name, "addon") {
			metadata[config.OAMComponentCategoryMetadataKey] = "addon"
		}

		oamRDP = append(oamRDP, adapter.OAMRegistrantDefinitionPath{
			OAMDefintionPath: pathSet.oamDefinitionPath,
			OAMRefSchemaPath: pathSet.jsonSchemaPath,
			Host:             host,
			Metadata:         metadata,
		})
	}

	return adapter.
		NewOAMRegistrant(oamRDP, fmt.Sprintf("%s/api/oam/workload", runtime)).
		Register()
}

// RegisterTraits will register all of the trait definitions
// present in the path oam/traits
//
// Registeration process will send POST request to $runtime/api/oam/trait
func RegisterTraits(runtime, host string) error {
	oamRDP := []adapter.OAMRegistrantDefinitionPath{}

	pathSets, err := load(traitPath)
	if err != nil {
		return err
	}

	for _, pathSet := range pathSets {
		metadata := map[string]string{
			config.OAMAdapterNameMetadataKey: config.KumaOperation,
		}

		oamRDP = append(oamRDP, adapter.OAMRegistrantDefinitionPath{
			OAMDefintionPath: pathSet.oamDefinitionPath,
			OAMRefSchemaPath: pathSet.jsonSchemaPath,
			Host:             host,
			Metadata:         metadata,
		})
	}

	return adapter.
		NewOAMRegistrant(oamRDP, fmt.Sprintf("%s/api/oam/trait", runtime)).
		Register()
}

func load(basePath string) ([]schemaDefinitionPathSet, error) {
	res := []schemaDefinitionPathSet{}

	if err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*_definition.json", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			nameWithPath := strings.TrimSuffix(path, "_definition.json")

			res = append(res, schemaDefinitionPathSet{
				oamDefinitionPath: path,
				jsonSchemaPath:    fmt.Sprintf("%s.meshery.layer5io.schema.json", nameWithPath),
				name:              filepath.Base(nameWithPath),
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

// type DynamicComponentsConfig struct {
// 	TimeoutInMinutes time.Duration
// 	Url              string
// 	GenerationMethod string
// 	Config           manifests.Config
// 	Operation        string
// }

// func RegisterWorkLoadsDynamically(runtime, host string, dc *DynamicComponentsConfig) error {
// 	var comp *manifests.Component
// 	var err error
// 	switch dc.GenerationMethod {
// 	case MANIFESTS:
// 		comp, err = manifests.GetFromManifest(dc.Url, manifests.SERVICE_MESH, dc.Config)
// 	case HELM_CHARTS:
// 		comp, err = manifests.GetFromHelm(dc.Url, manifests.SERVICE_MESH, dc.Config)
// 	default:
// 		return err
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	for i, def := range comp.Definitions {
// 		var ord adapter.OAMRegistrantData
// 		ord.OAMRefSchema = comp.Schemas[i]

// 		//Marshalling the stringified json
// 		ord.Host = host
// 		definitionMap := map[string]interface{}{}
// 		if err := json.Unmarshal([]byte(def), &definitionMap); err != nil {
// 			return err
// 		}
// 		definitionMap["apiVersion"] = "core.oam.dev/v1alpha1"
// 		definitionMap["kind"] = "WorkloadDefinition"
// 		ord.OAMDefinition = definitionMap
// 		ord.Metadata = map[string]string{
// 			OAMAdapterNameMetadataKey: dc.Operation,
// 		}
// 		// send request to the register
// 		backoffOpt := backoff.NewExponentialBackOff()
// 		backoffOpt.MaxElapsedTime = 10 * dc.TimeoutInMinutes
// 		if err := backoff.Retry(func() error {
// 			contentByt, err := json.Marshal(ord)
// 			if err != nil {
// 				return backoff.Permanent(err)
// 			}
// 			content := bytes.NewReader(contentByt)
// 			// host here is given by the application itself and is trustworthy hence,
// 			// #nosec
// 			resp, err := http.Post(fmt.Sprintf("%s/api/oam/workload", runtime), "application/json", content)
// 			if err != nil {
// 				return err
// 			}
// 			if resp.StatusCode != http.StatusCreated &&
// 				resp.StatusCode != http.StatusOK &&
// 				resp.StatusCode != http.StatusAccepted {
// 				return fmt.Errorf(
// 					"register process failed, host returned status: %s with status code %d",
// 					resp.Status,
// 					resp.StatusCode,
// 				)
// 			}

// 			return nil
// 		}, backoffOpt); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
