module github.com/layer5io/meshery-kuma

go 1.13

replace (
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
//github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
//golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
)

require (
	github.com/layer5io/meshery-adapter-library v0.5.3
	github.com/layer5io/meshkit v0.5.8
	github.com/layer5io/service-mesh-performance v0.3.4
	gopkg.in/yaml.v2 v2.4.0 // direct
	k8s.io/apimachinery v0.21.0 // direct
)
