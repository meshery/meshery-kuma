package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-kuma/kuma"
	"github.com/layer5io/meshery-kuma/kuma/oam"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultVersion string
var DefaultGenerationURL string
var DefaultGenerationMethod string
var WorkloadPath string
var MeshModelPath string

var Meshmodelmetadata = make(map[string]interface{})

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category:    "Orchestration & Management",
	SubCategory: "Service Mesh",
	Metadata:    Meshmodelmetadata,
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_KUMA)],
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}
func getLatestValidAppVersionAndChartVersion() (string, string, error) {
	release, err := utils.GetLatestReleaseTagsSorted("kumahq", "kuma")
	if err != nil {
		return "", "", kuma.ErrGetLatestRelease(err)
	}
	//loops through latest 10 app versions until it finds one which is available in helm chart's index.yaml
	for i := range release {
		if chartVersion, err := kubernetes.HelmAppVersionToChartVersion("https://kumahq.github.io/charts", "kuma", release[len(release)-i-1]); err == nil {
			return release[len(release)-i-1], chartVersion, nil
		}
	}
	return "", "", kuma.ErrGetLatestRelease(err)
}

func init() {
	//Initialize Metadata including logo svgs
	f, _ := os.Open("./build/meshmodel_metadata.json")
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()
	byt, _ := io.ReadAll(f)

	_ = json.Unmarshal(byt, &Meshmodelmetadata)
	wd, _ := os.Getwd()

	var chartVersion string
	DefaultVersion, chartVersion, _ = getLatestValidAppVersionAndChartVersion()
	DefaultGenerationURL = "https://github.com/kumahq/charts/releases/download/kuma-" + chartVersion + "/kuma-" + chartVersion + ".tgz"
	DefaultGenerationMethod = adapter.HelmCHARTS
	WorkloadPath = oam.WorkloadPath
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
}
