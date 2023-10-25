module github.com/layer5io/meshery-kuma

go 1.21

replace (

	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v0.11.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v0.11.0
)


require (
	github.com/google/uuid v1.3.1
	github.com/layer5io/meshery-adapter-library v0.6.7
	github.com/layer5io/meshkit v0.6.77
	github.com/layer5io/service-mesh-performance v0.3.4
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.26.1
)
