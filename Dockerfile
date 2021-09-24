FROM golang:1.13 as builder

ARG VERSION
ARG GIT_COMMITSHA
WORKDIR /build
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN GOPROXY=https://proxy.golang.org,direct go mod download
# Copy the go source
COPY main.go main.go
COPY internal/ internal/
COPY kuma/ kuma/
# Build
RUN GOPROXY=https://proxy.golang.org,direct CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -a -o meshery-kuma main.go

FROM alpine:3.14 as jsonschema-util
RUN apk add --no-cache curl
WORKDIR /
RUN curl -LO https://github.com/layer5io/kubeopenapi-jsonschema/releases/download/v0.1.0/kubeopenapi-jsonschema
RUN chmod +x /kubeopenapi-jsonschema

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/nodejs:14
WORKDIR /
ENV DISTRO="debian"
ENV GOARCH="amd64"
ENV SERVICE_ADDR="meshery-kuma"
ENV MESHERY_SERVER="http://meshery:9081"
COPY templates/ ./templates
COPY --from=builder /build/meshery-kuma .
COPY --from=jsonschema-util /kubeopenapi-jsonschema /root/.meshery/bin/kubeopenapi-jsonschema
ENTRYPOINT ["/meshery-kuma"]
