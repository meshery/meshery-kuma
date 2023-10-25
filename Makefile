
include build/Makefile.core.mk
include build/Makefile.show-help.mk

#-----------------------------------------------------------------------------
# Environment Setup
#-----------------------------------------------------------------------------
BUILDER=buildx-multi-arch
ADAPTER=kuma



#-----------------------------------------------------------------------------
# Docker-based Builds
#-----------------------------------------------------------------------------
.PHONY: docker docker-run lint error test run run-force-dynamic-reg


## Lint check Golang
lint:
	golangci-lint run -c .golangci.yml -v ./...

tidy:
	go mod tidy

build:
	go build -o ./bin/$(ADAPTER) ./main.go

## Build Adapter container image with "edge-latest" tag
docker:
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery-$(ADAPTER):$(RELEASE_CHANNEL)-latest .

## Run Adapter container with "edge-latest" tag
docker-run:
	(docker rm -f meshery-$(ADAPTER)) || true
	docker run --name meshery-$(ADAPTER) -d \
	-p 10000:10000 \
	-e DEBUG=true \
	layer5/meshery-$(ADAPTER):$(RELEASE_CHANNEL)-latest

run:
	go run main.go

## Build and run Adapter locally; force component registration
run-force-dynamic-reg:
	FORCE_DYNAMIC_REG=true DEBUG=true GOSUMDB=off go run main.go

## Run Meshery Error utility
error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers

## Run Golang tests
test:
	export CURRENTCONTEXT="$(kubectl config current-context)" 
	echo "current-context:" ${CURRENTCONTEXT} 
	export KUBECONFIG="${HOME}/.kube/config"
	echo "environment-kubeconfig:" ${KUBECONFIG}
	GOPROXY=direct GOSUMDB=off GO111MODULE=on go test -v ./...
