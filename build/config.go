package build

import (
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

//NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_KUMA)],
		MeshVersion: version,
		Filter: manifests.CrdFilter{
			RootFilter:    []string{"$[?(@.kind==\"CustomResourceDefinition\")]"},
			NameFilter:    []string{"$..[\"spec\"][\"names\"][\"kind\"]"},
			VersionFilter: []string{"$[0]..spec.versions[0]"},
			GroupFilter:   []string{"$[0]..spec"},
			SpecFilter:    []string{"$[0]..openAPIV3Schema.properties.spec"},
			ItrFilter:     []string{"$[?(@.spec.names.kind"},
			ItrSpecFilter: []string{"$[?(@.spec.names.kind"},
			VField:        "name",
			GField:        "group",
		},
	}
}
func getLatestValidAppVersionAndChartVersion() (string, string, error) {
	release, err := utils.GetLatestReleaseTagsSorted("kumahq", "kuma")
	if err != nil {
		return "", "", kuma.ErrGetLatestRelease(err)
	}
	//loops through latest 10 app versions untill it finds one which is available in helm chart's index.yaml
	for i := range release {
		if chartVersion, err := kubernetes.HelmAppVersionToChartVersion("https://kumahq.github.io/charts", "kuma", release[len(release)-i-1]); err == nil {
			return release[len(release)-i-1], chartVersion, nil
		}

	}
	return "", "", kuma.ErrGetLatestRelease(err)
}

func init() {
	var chartVersion string
	DefaultVersion, chartVersion, _ = getLatestValidAppVersionAndChartVersion()
	DefaultGenerationURL = "https://github.com/kumahq/charts/releases/download/kuma-" + chartVersion + "/kuma-" + chartVersion + ".tgz"
	DefaultGenerationMethod = adapter.HelmCHARTS
	WorkloadPath = oam.WorkloadPath
}
