module github.com/layer5io/meshery-kuma

go 1.13

replace (
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
)

require (
	github.com/layer5io/meshery-adapter-library v0.1.13
	github.com/layer5io/meshkit v0.2.6
	github.com/layer5io/service-mesh-performance v0.3.3
	helm.sh/helm/v3 v3.3.4 // indirect
)
