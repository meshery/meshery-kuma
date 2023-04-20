package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-kuma/kuma"
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

var meshmodelmetadata = map[string]interface{}{
	"Primary Color":   "#291953",
	"Secondary Color": "#6942c9",
	"Shape":           "circle",
	"Logo URL":        "https://github.com/cncf/artwork/blob/master/projects/kuma/icon/white/kuma-icon-white.svg",
	"SVG_Color":       "<svg id=\"Layer_1\" data-name=\"Layer 1\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 1173.18 1173.18\"><defs><style>.cls-1{fill:#291953;}.cls-2{fill:none;}</style></defs><g id=\"Layer_2\" data-name=\"Layer 2\"><g id=\"Layer_1-2\" data-name=\"Layer 1-2\"><path class=\"cls-1\" d=\"M887.62,422.54a6.21,6.21,0,0,1,1-5.9c24.85-31.37,47.4-67.46,47.4-95.14C936,260,900.91,210,824.51,210c-37.85,0-65.61,12.3-83.86,32.11a6.39,6.39,0,0,1-6.68,1.8,570.26,570.26,0,0,0-89.24-21.12,6.24,6.24,0,0,0-7,5.35,6.14,6.14,0,0,0,.16,2.45c6.31,23.66,44.2,174,74.71,288.44,18.45,69.26-29.36,137.3-101,137.09H567.19c-72.42,0-116.38-68.28-99.69-136.35,28.17-115,66.76-264.17,73-288.77a6.19,6.19,0,0,0-4.37-7.59,6,6,0,0,0-2.39-.16,486.69,486.69,0,0,0-103.38,23.66,6.37,6.37,0,0,1-7-1.93c-18.24-21.45-46.7-34.86-86.11-34.86-76.4,0-111.5,49.91-111.5,111.5,0,32.28,30.67,76,59.87,110.31a6.36,6.36,0,0,1,1.15,6.07l-49.7,144.35a1.14,1.14,0,0,0,0,.45c-1.31,5-20.51,90.22,125.32,225.79C406,849.23,558,995.66,585.35,1021.83a6.16,6.16,0,0,0,8.49,0c28.09-26.13,185.77-172.48,229.65-213.24,157.55-146.93,120-226.24,120-226.24Z\"/><path class=\"cls-1\" d=\"M619.23,560.53H559.85a17.8,17.8,0,0,1-17.8-17.79v-.09l-7.38-73.11a17.8,17.8,0,0,1,17.8-17.8h73.85a17.8,17.8,0,0,1,17.84,17.76v0l-7.09,73.11a17.8,17.8,0,0,1-17.72,17.88Z\"/><rect class=\"cls-2\" width=\"1173.18\" height=\"1173.18\"/></g></g></svg>",
	"SVG_White":       "<svg id=\"Layer_1\" data-name=\"Layer 1\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 1173.18 1173.18\"><defs><style>.cls-1{fill:#fff;}.cls-2{fill:none;}</style></defs><g id=\"Layer_2\" data-name=\"Layer 2\"><g id=\"Layer_1-2\" data-name=\"Layer 1-2\"><path class=\"cls-1\" d=\"M887.62,422.54a6.21,6.21,0,0,1,1-5.9c24.85-31.37,47.4-67.46,47.4-95.14C936,260,900.91,210,824.51,210c-37.85,0-65.61,12.3-83.86,32.11a6.39,6.39,0,0,1-6.68,1.8,570.26,570.26,0,0,0-89.24-21.12,6.24,6.24,0,0,0-7,5.35,6.14,6.14,0,0,0,.16,2.45c6.31,23.66,44.2,174,74.71,288.44,18.45,69.26-29.36,137.3-101,137.09H567.19c-72.42,0-116.38-68.28-99.69-136.35,28.17-115,66.76-264.17,73-288.77a6.19,6.19,0,0,0-4.37-7.59,6,6,0,0,0-2.39-.16,486.69,486.69,0,0,0-103.38,23.66,6.37,6.37,0,0,1-7-1.93c-18.24-21.45-46.7-34.86-86.11-34.86-76.4,0-111.5,49.91-111.5,111.5,0,32.28,30.67,76,59.87,110.31a6.36,6.36,0,0,1,1.15,6.07l-49.7,144.35a1.14,1.14,0,0,0,0,.45c-1.31,5-20.51,90.22,125.32,225.79C406,849.23,558,995.66,585.35,1021.83a6.16,6.16,0,0,0,8.49,0c28.09-26.13,185.77-172.48,229.65-213.24,157.55-146.93,120-226.24,120-226.24Z\"/><path class=\"cls-1\" d=\"M619.23,560.53H559.85a17.8,17.8,0,0,1-17.8-17.79v-.09l-7.38-73.11a17.8,17.8,0,0,1,17.8-17.8h73.85a17.8,17.8,0,0,1,17.84,17.76v0l-7.09,73.11a17.8,17.8,0,0,1-17.72,17.88Z\"/><rect class=\"cls-2\" width=\"1173.18\" height=\"1173.18\"/></g></g></svg>",
}

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category:    "Cloud Native Network",
	SubCategory: "Service Mesh",
	Metadata:    meshmodelmetadata,
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
	wd, _ := os.Getwd()

	var chartVersion string
	DefaultVersion, chartVersion, _ = getLatestValidAppVersionAndChartVersion()
	DefaultGenerationURL = "https://github.com/kumahq/charts/releases/download/kuma-" + chartVersion + "/kuma-" + chartVersion + ".tgz"
	DefaultGenerationMethod = adapter.HelmCHARTS
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
}
